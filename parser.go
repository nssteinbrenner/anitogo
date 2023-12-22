package anitogo

import (
	"regexp"
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

func (p *parser) parse() {
	p.searchForKeywords()
	p.searchForIsolatedNumbers()
	if p.tokenizer.options.ParseEpisodeNumber {
		p.searchForEpisodeNumber()
	}
	p.searchForAnimeTitle()
	if p.tokenizer.options.ParseReleaseGroup && !p.tokenizer.elements.contains(elementCategoryReleaseGroup) {
		p.searchForReleaseGroup()
	}
	if p.tokenizer.options.ParseEpisodeTitle && p.tokenizer.elements.contains(elementCategoryEpisodeNumber) {
		p.searchForEpisodeTitle()
	}
	p.validateElements()
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
			cat = kd.category
			if !p.tokenizer.options.ParseReleaseGroup && cat == elementCategoryReleaseGroup {
				continue
			}
			if !cat.isSearchable() || !kd.options.searchable {
				continue
			}
			if cat.isSingular() && p.tokenizer.elements.contains(cat) {
				continue
			}

			if cat == elementCategoryAnimeSeasonPrefix {
				p.checkAnimeSeasonKeyword(tkn)
				continue
			} else if cat == elementCategoryEpisodePrefix {
				if kd.options.valid {
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
			if kd.empty() || kd.options.identifiable {
				tkn.Category = tokenCategoryIdentifier
			}
		}
	}
}

func (p *parser) searchForIsolatedNumbers() {
	for _, tkn := range p.tokenizer.tokens.getListFlag(tokenFlagsUnknown) {
		if !isNumeric(tkn.Content) {
			continue
		}
		isolated := p.tokenizer.tokens.isTokenIsolated(*tkn)
		if !isolated {
			continue
		}

		n, _ := strconv.Atoi(tkn.Content)

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
}

func (p *parser) searchForEpisodeNumber() {
	tkns := p.tokenizer.tokens.getListFlag(tokenFlagsUnknown)
	if len(tkns) == 0 {
		return
	}

	p.tokenizer.elements.setCheckAltNumber(p.tokenizer.elements.contains(elementCategoryEpisodeNumber))

	match := p.searchForEpisodePatterns(tkns)
	if match {
		return
	}

	if p.tokenizer.elements.contains(elementCategoryEpisodeNumber) {
		return
	}

	var numericTokens tokens
	for _, v := range tkns {
		if isNumeric(v.Content) {
			numericTokens = append(numericTokens, v)
		}
	}

	if len(numericTokens) == 0 {
		return
	}

	if p.searchForEquivalentNumbers(numericTokens) {
		return
	}

	if p.searchForSeparatedNumbers(numericTokens) {
		return
	}

	if p.searchForIsolatedNumbersTokens(numericTokens) {
		return
	}

	if p.searchForLastNumber(numericTokens) {
		return
	}
}

func (p *parser) searchForAnimeTitle() {
	enclosedTitle := false

	tokenBegin, found := p.tokenizer.tokens.find(tokenFlagsNotEnclosed | tokenFlagsUnknown)
	if !found {
		enclosedTitle = true
		tokenBegin, found = p.tokenizer.tokens.get(0)
		skippedPreviousGroup := false
		for found {
			tokenBegin, found = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsUnknown)
			if !found {
				break
			}
			if isMostlyLatinString(tokenBegin.Content) {
				if skippedPreviousGroup {
					break
				}
			}
			tokenBegin, found = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsBracket)
			skippedPreviousGroup = true
		}
	}
	if tokenBegin.empty() {
		return
	}

	targetFlag := tokenFlagsNone
	if enclosedTitle {
		targetFlag = tokenFlagsBracket
	}
	tokenEnd, _ := p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsIdentifier|targetFlag)
	if !enclosedTitle {
		lastBracket := tokenEnd
		bracketOpen := false
		tknList := p.tokenizer.tokens.getList(tokenFlagsBracket, tokenBegin, tokenEnd)
		for _, tkn := range tknList {
			lastBracket = tkn
			bracketOpen = !bracketOpen
		}
		if bracketOpen {
			tokenEnd = lastBracket
		}
	}

	if !enclosedTitle {
		tkn, found := p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsNotDelimiter)
		if !found {
			return
		}
		for tkn.Category == tokenCategoryBracket && tkn.Content != ")" {
			tkn, found = p.tokenizer.tokens.findPrevious(*tkn, tokenFlagsBracket)
			if found {
				if !tkn.empty() {
					tokenEnd = tkn
					tkn, _ = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsNotDelimiter)
				}
			}
		}
	}

	tokenEnd, _ = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsValid)
	p.buildElement(elementCategoryAnimeTitle, tokenBegin, tokenEnd, false)
}

func (p *parser) searchForReleaseGroup() {
	tokenEnd := &token{}
	tokenBegin := &token{}
	previousToken := &token{}
	for {
		if !tokenEnd.empty() {
			tokenBegin, _ = p.tokenizer.tokens.findNext(*tokenEnd, tokenFlagsEnclosed|tokenFlagsUnknown)
		} else {
			tokenBegin, _ = p.tokenizer.tokens.find(tokenFlagsEnclosed | tokenFlagsUnknown)
		}
		if tokenBegin.empty() {
			return
		}
		tokenEnd, _ = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsBracket|tokenFlagsIdentifier)
		if tokenEnd.empty() {
			return
		}
		if tokenEnd.Category != tokenCategoryBracket {
			continue
		}
		previousToken, _ = p.tokenizer.tokens.findPrevious(*tokenBegin, tokenFlagsNotDelimiter)
		if !previousToken.empty() && previousToken.Category != tokenCategoryBracket {
			continue
		}

		tokenEnd, _ = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsValid)
		p.buildElement(elementCategoryReleaseGroup, tokenBegin, tokenEnd, true)
		return
	}
}

func (p *parser) searchForEpisodeTitle() {
	tokenEnd := &token{}
	tokenBegin := &token{}
	for {
		if !tokenEnd.empty() {
			tokenBegin, _ = p.tokenizer.tokens.findNext(*tokenEnd, tokenFlagsNotEnclosed|tokenFlagsUnknown)
		} else {
			tokenBegin, _ = p.tokenizer.tokens.find(tokenFlagsNotEnclosed | tokenFlagsUnknown)
		}
		if tokenBegin.empty() {
			return
		}
		tokenEnd, _ = p.tokenizer.tokens.findNext(*tokenBegin, tokenFlagsBracket|tokenFlagsIdentifier)
		if tokenEnd.empty() {
			tokenEnd, _ = p.tokenizer.tokens.get(len(*p.tokenizer.tokens) - 1)
		}
		dist := p.tokenizer.tokens.distance(tokenBegin, tokenEnd)
		if dist >= 0 && dist <= 2 && isDashCharacter(tokenBegin.Content) {
			continue
		}

		if !tokenEnd.empty() && tokenEnd.Category == tokenCategoryBracket {
			tokenEnd, _ = p.tokenizer.tokens.findPrevious(*tokenEnd, tokenFlagsValid)
		}
		p.buildElement(elementCategoryEpisodeTitle, tokenBegin, tokenEnd, false)
		return
	}
}

func (p *parser) validateElements() {
	if p.tokenizer.elements.contains(elementCategoryEpisodeTitle) {
		episodeTitle := p.tokenizer.elements.get(elementCategoryEpisodeTitle)[0]
		re := regexp.MustCompile("^~\\s(\\d{1,2})$")
		match := re.FindStringSubmatch(episodeTitle)
		if match != nil {
			p.tokenizer.elements.erase(elementCategoryEpisodeTitle)
			p.tokenizer.elements.insert(elementCategoryEpisodeNumber, match[1])
		}
	}
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
