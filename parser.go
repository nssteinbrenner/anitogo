package anitogo

import (
	"strconv"
	"strings"
)

type parser struct {
	tokenizer *tokenizer
}

func newParser(tkz *tokenizer) *parser {
	psr := parser{
		tokenizer: tkz,
	}
	return &psr
}

func (p *parser) parse() error {
	p.searchForKeywords()
	err := p.searchForIsolatedNumbers()
	if err != nil {
		return err
	}
	if p.tokenizer.options.ParseEpisodeNumber {
		err = p.searchForEpisodeNumber()
		if err != nil {
			return err
		}
	}
	err = p.searchForAnimeTitle()
	if err != nil {
		return err
	}
	if p.tokenizer.options.ParseReleaseGroup && !p.tokenizer.elements.contains(elementCategoryReleaseGroup) {
		err = p.searchForReleaseGroup()
		if err != nil {
			return err
		}
	}
	if p.tokenizer.options.ParseEpisodeTitle && p.tokenizer.elements.contains(elementCategoryEpisodeNumber) {
		err = p.searchForEpisodeTitle()
		if err != nil {
			return err
		}
	}
	p.validateElements()
	return nil
}

func (p *parser) searchForKeywords() {
	for _, tkn := range p.tokenizer.tokens.getListFlag(tokenFlagsUnknown) {
		w := tkn.Content
		w = strings.Trim(w, " -")

		if w == "" {
			continue
		}
		if len(w) != 8 && isNumeric(w) {
			continue
		}

		cat := elementCategoryUnknown
		kd, found := p.tokenizer.keywordManager.findWithoutCategory(p.tokenizer.keywordManager.normalize(w))
		if found {
			cat = kd.Category
			if !p.tokenizer.options.ParseReleaseGroup && cat == elementCategoryReleaseGroup {
				continue
			}
			if !cat.isSearchable() || !kd.Options.Searchable {
				continue
			}
			if cat.isSingular() && p.tokenizer.elements.contains(cat) {
				continue
			}

			if cat == elementCategoryAnimeSeasonPrefix {
				p.checkAnimeSeasonKeyword(tkn)
				continue
			} else if cat == elementCategoryEpisodePrefix {
				if kd.Options.Valid {
					p.checkExtentKeyword(elementCategoryEpisodeNumber, tkn)
				}
				continue
			} else if cat == elementCategoryReleaseVersion {
				w = w[1:]
			} else if cat == elementCategoryVolumePrefix {
				p.checkExtentKeyword(elementCategoryVolumeNumber, tkn)
				continue
			}
		} else {
			if !p.tokenizer.elements.contains(elementCategoryFileChecksum) && isCRC32(w) {
				cat = elementCategoryFileChecksum
			} else if !p.tokenizer.elements.contains(elementCategoryVideoResolution) && isResolution(w) {
				cat = elementCategoryVideoResolution
			}
		}

		if cat != elementCategoryUnknown {
			p.tokenizer.elements.insert(cat, w)
			if kd.empty() || kd.Options.Identifiable {
				tkn.Category = tokenCategoryIdentifier
			}
		}
	}
}

func (p *parser) searchForIsolatedNumbers() error {
	for _, tkn := range p.tokenizer.tokens.getListFlag(tokenFlagsUnknown) {
		if !isNumeric(tkn.Content) {
			continue
		}
		isolated, err := p.tokenizer.tokens.isTokenIsolated(*tkn)
		if err != nil {
			return err
		}
		if !isolated {
			continue
		}

		n, err := strconv.Atoi(tkn.Content)
		if err != nil {
			continue
		}

		if n >= animeYearMin && n <= animeYearMax {
			if !p.tokenizer.elements.contains(elementCategoryAnimeYear) {
				p.tokenizer.elements.insert(elementCategoryAnimeYear, tkn.Content)
				tkn.Category = tokenCategoryIdentifier
				continue
			}
		}

		if n == 480 || n == 720 || n == 1080 {
			if !p.tokenizer.elements.contains(elementCategoryVideoResolution) {
				p.tokenizer.elements.insert(elementCategoryVideoResolution, tkn.Content)
				tkn.Category = tokenCategoryIdentifier
				continue
			}
		}
	}
	return nil
}

func (p *parser) searchForEpisodeNumber() error {
	tkns := p.tokenizer.tokens.getListFlag(tokenFlagsUnknown)
	if len(tkns) == 0 {
		return nil
	}

	p.tokenizer.elements.setCheckAltNumber(p.tokenizer.elements.contains(elementCategoryEpisodeNumber))

	match, err := p.searchForEpisodePatterns(tkns)
	if err != nil {
		return err
	}
	if match {
		return nil
	}

	if p.tokenizer.elements.contains(elementCategoryEpisodeNumber) {
		return nil
	}

	var numericTokens tokens
	for _, v := range tkns {
		if isNumeric(v.Content) {
			numericTokens = append(numericTokens, v)
		}
	}

	if len(numericTokens) == 0 {
		return nil
	}

	match, err = p.searchForEquivalentNumbers(numericTokens)
	if err != nil {
		return err
	}
	if match {
		return nil
	}

	match, err = p.searchForSeparatedNumbers(numericTokens)
	if err != nil {
		return err
	}
	if match {
		return nil
	}

	match, err = p.searchForIsolatedNumbersTokens(numericTokens)
	if err != nil {
		return err
	}
	if match {
		return nil
	}

	match, err = p.searchForLastNumber(numericTokens)
	if err != nil {
		return err
	}
	if match {
		return nil
	}
	return nil
}

func (p *parser) searchForAnimeTitle() error {
	var err error
	enclosedTitle := false

	tokenBegin, found := p.tokenizer.tokens.find(tokenFlagsNotEnclosed | tokenFlagsUnknown)
	if !found {
		enclosedTitle = true
		tokenBegin, found = p.tokenizer.tokens.get(0)
		skippedPreviousGroup := false
		for found {
			tokenBegin, found, err = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsUnknown)
			if err != nil {
				return err
			}
			if !found {
				break
			}
			if isMostlyLatinString(tokenBegin.Content) {
				if skippedPreviousGroup {
					break
				}
			}
			tokenBegin, found, err = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsBracket)
			if err != nil {
				return err
			}
			skippedPreviousGroup = true
		}
	}
	if tokenBegin.empty() {
		return nil
	}

	targetFlag := tokenFlagsNone
	if enclosedTitle {
		targetFlag = tokenFlagsBracket
	}
	tokenEnd, _, err := p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsIdentifier|targetFlag)
	if err != nil {
		return err
	}
	if !enclosedTitle {
		lastBracket := tokenEnd
		bracketOpen := false
		tknList, err := p.tokenizer.tokens.getList(tokenFlagsBracket, tokenBegin, tokenEnd)
		if err != nil {
			return err
		}
		for _, tkn := range tknList {
			lastBracket = tkn
			bracketOpen = !bracketOpen
		}
		if bracketOpen {
			tokenEnd = lastBracket
		}
	}

	if !enclosedTitle {
		tkn, _, err := p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsNotDelimiter)
		if err == nil {
			for tkn.Category == tokenCategoryBracket && tkn.Content != ")" {
				tkn, _, err = p.tokenizer.tokens.findPrevious(*tkn, tokenFlagsBracket)
				if err != nil {
					return err
				}
				if !tkn.empty() {
					tokenEnd = tkn
					tkn, _, err = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsNotDelimiter)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	tokenEnd, _, _ = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsValid)
	p.buildElement(elementCategoryAnimeTitle, tokenBegin, tokenEnd, false)
	return nil
}

func (p *parser) searchForReleaseGroup() error {
	var err error
	tokenEnd := &token{}
	tokenBegin := &token{}
	previousToken := &token{}
	for {
		if !tokenEnd.empty() {
			tokenBegin, _, err = p.tokenizer.tokens.findNext(*tokenEnd, tokenFlagsEnclosed|tokenFlagsUnknown)
			if err != nil {
				return err
			}
		} else {
			tokenBegin, _ = p.tokenizer.tokens.find(tokenFlagsEnclosed | tokenFlagsUnknown)
		}
		if tokenBegin.empty() {
			return nil
		}
		tokenEnd, _, err = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsBracket|tokenFlagsIdentifier)
		if err != nil {
			return err
		}
		if tokenEnd.empty() {
			return nil
		}
		if tokenEnd.Category != tokenCategoryBracket {
			continue
		}
		previousToken, _, err = p.tokenizer.tokens.findPrevious(*tokenBegin, tokenFlagsNotDelimiter)
		if err != nil {
			return err
		}
		if !previousToken.empty() && previousToken.Category != tokenCategoryBracket {
			continue
		}

		tokenEnd, _, err = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsValid)
		if err != nil {
			return err
		}
		p.buildElement(elementCategoryReleaseGroup, tokenBegin, tokenEnd, true)
		return nil
	}
}

func (p *parser) searchForEpisodeTitle() error {
	var err error
	tokenEnd := &token{}
	tokenBegin := &token{}
	for {
		if !tokenEnd.empty() {
			tokenBegin, _, err = p.tokenizer.tokens.findNext(*tokenEnd, tokenFlagsNotEnclosed|tokenFlagsUnknown)
			if err != nil {
				return err
			}
		} else {
			tokenBegin, _ = p.tokenizer.tokens.find(tokenFlagsNotEnclosed | tokenFlagsUnknown)
		}
		if tokenBegin.empty() {
			return nil
		}
		tokenEnd, _, err = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsBracket|tokenFlagsIdentifier)
		if err != nil {
			return err
		}
		if tokenEnd.empty() {
			tokenEnd, _ = p.tokenizer.tokens.get(len(*p.tokenizer.tokens) - 1)
		}
		dist, err := p.tokenizer.tokens.distance(tokenBegin, tokenEnd)
		if err != nil {
			return err
		}
		if dist <= 2 && isDashCharacter(tokenBegin.Content) {
			continue
		}

		if !tokenEnd.empty() && tokenEnd.Category == tokenCategoryBracket {
			tokenEnd, _, err = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsValid)
			if err != nil {
				return err
			}
		}
		p.buildElement(elementCategoryEpisodeTitle, tokenBegin, tokenEnd, false)
		return nil
	}
}

func (p *parser) validateElements() {
	if p.tokenizer.elements.contains(elementCategoryAnimeType) && p.tokenizer.elements.contains(elementCategoryEpisodeTitle) {
		episodeTitle := p.tokenizer.elements.get(elementCategoryEpisodeTitle)[0]
		animeTypeList := p.tokenizer.elements.get(elementCategoryAnimeType)
		for _, animeType := range animeTypeList {
			if animeType == episodeTitle {
				p.tokenizer.elements.erase(elementCategoryEpisodeTitle)
			} else if strings.Contains(episodeTitle, animeType) {
				normAnimeType := p.tokenizer.keywordManager.normalize(animeType)
				_, found := p.tokenizer.keywordManager.find(normAnimeType, elementCategoryAnimeType)
				if found {
					p.tokenizer.elements.remove(elementCategoryAnimeType, animeType)
				}
				continue
			}
		}
	}
}
