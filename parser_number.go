package anitogo

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	animeYearMin     = 1900
	animeYearMax     = 2050
	episodeNumberMax = 1899
	volumeNumberMax  = 20
)

func (p *parser) checkExtentKeyword(cat elementCategory, tkn *token) bool {
	nextToken, _ := p.tokenizer.tokens.findNext(*tkn, tokenFlagsNotDelimiter)

	if nextToken.Category == tokenCategoryUnknown {
		if !nextToken.empty() && findNumberInString(nextToken.Content) > -1 {
			if cat == elementCategoryEpisodeNumber {
				match := p.matchEpisodePattern(nextToken.Content, nextToken)
				if !match {
					p.setEpisodeNumber(nextToken.Content, nextToken, false)
				}
			} else if cat == elementCategoryVolumeNumber {
				if !p.matchVolumePattern(nextToken.Content, nextToken) {
					p.setVolumeNumber(nextToken.Content, nextToken, false)
				}
			} else {
				return false
			}
			tkn.Category = tokenCategoryIdentifier
			return true
		}
	}

	return false
}

func (p *parser) searchForEpisodePatterns(tkns tokens) bool {
	for _, tkn := range tkns {
		numericFront := isNumeric(string(tkn.Content[0]))

		if !numericFront {
			if p.numberComesAfterPrefix(elementCategoryEpisodePrefix, tkn) {
				return true
			}
			if p.numberComesAfterPrefix(elementCategoryVolumePrefix, tkn) {
				continue
			}
			if p.numberComesAfterPrefix(elementCategoryAnimeSeasonPrefix, tkn) {
				continue
			}
		} else {
			if p.numberComesBeforeAnotherNumber(tkn) {
				return true
			}
		}
		if p.matchEpisodePattern(tkn.Content, tkn) {
			return true
		}
	}
	return false
}

func (p *parser) numberComesAfterPrefix(cat elementCategory, tkn *token) bool {
	numberBegin := findNumberInString(tkn.Content)
	if numberBegin == -1 {
		return false
	}
	prefix := tkn.Content[:numberBegin]

	_, found := p.tokenizer.keywordManager.find(p.tokenizer.keywordManager.normalize(prefix), cat)
	if found {
		number := tkn.Content[numberBegin:]
		if cat == elementCategoryEpisodePrefix {
			if p.matchEpisodePattern(number, tkn) {
				return true
			}
			return p.setEpisodeNumber(number, tkn, false)
		}
		if cat == elementCategoryVolumePrefix {
			if p.matchVolumePattern(number, tkn) {
				return true
			}
			return p.setVolumeNumber(number, tkn, false)
		}
		if cat == elementCategoryAnimeSeasonPrefix {
			return p.setSeasonNumber(number, tkn)
		}
	}
	return false
}

func (p *parser) numberComesBeforeAnotherNumber(tkn *token) bool {
	separatorToken, found := p.tokenizer.tokens.findNext(*tkn, tokenFlagsNotDelimiter)

	if found {
		separator := separatorToken.Content
		if separator == "&" || separator == "of" {
			otherToken, found := p.tokenizer.tokens.findNext(*separatorToken, tokenFlagsNotDelimiter)
			if found && isNumeric(otherToken.Content) {
				p.setEpisodeNumber(tkn.Content, tkn, false)
				if separator == "&" {
					p.setEpisodeNumber(otherToken.Content, tkn, false)
				}
				separatorToken.Category = tokenCategoryIdentifier
				otherToken.Category = tokenCategoryIdentifier
				return true
			}
		}
	}
	return false
}

func (p *parser) searchForEquivalentNumbers(tkns tokens) bool {
	for _, tkn := range tkns {
		if p.tokenizer.tokens.isTokenIsolated(*tkn) || !isValidEpisodeNumber(tkn.Content) {
			return false
		}

		nextToken, found := p.tokenizer.tokens.findNext(*tkn, tokenFlagsNotDelimiter)
		if nextToken.empty() || nextToken.Category != tokenCategoryBracket || !found {
			continue
		}
		nextToken, found = p.tokenizer.tokens.findNext(*nextToken, tokenFlagsEnclosed|tokenFlagsNotDelimiter)
		if found {
			if nextToken.Category != tokenCategoryUnknown {
				continue
			}
		}

		if !p.tokenizer.tokens.isTokenIsolated(*nextToken) || !isNumeric(nextToken.Content) || !isValidEpisodeNumber(nextToken.Content) {
			continue
		}

		i, _ := strconv.Atoi(nextToken.Content)
		j, _ := strconv.Atoi(tkn.Content)

		episode := nextToken
		altEpisode := tkn
		if i > j {
			episode = tkn
			altEpisode = nextToken
		}
		p.setEpisodeNumber(episode.Content, episode, false)
		p.setAlternativeEpisodeNumber(altEpisode.Content, altEpisode)
		return true
	}

	return false
}

func (p *parser) searchForSeparatedNumbers(tkns tokens) bool {
	for _, tkn := range tkns {
		previousToken, found := p.tokenizer.tokens.findPrevious(*tkn, tokenFlagsNotDelimiter)
		if !found {
			return false
		}

		if previousToken.Category == tokenCategoryUnknown && isDashCharacter(previousToken.Content) {
			if p.setEpisodeNumber(tkn.Content, tkn, true) {
				previousToken.Category = tokenCategoryIdentifier
				return true
			}
		}
	}
	return false
}

func (p *parser) searchForIsolatedNumbersTokens(tkns tokens) bool {
	for _, tkn := range tkns {
		if !tkn.Enclosed || !p.tokenizer.tokens.isTokenIsolated(*tkn) {
			continue
		}
		if p.setEpisodeNumber(tkn.Content, tkn, true) {
			return true
		}
	}
	return false
}

func (p *parser) searchForLastNumber(tkns tokens) bool {
	for _, tkn := range tkns {
		tokenIndex := p.tokenizer.tokens.getIndex(*tkn, 0)

		if tokenIndex == 0 {
			continue
		}

		if tkn.Enclosed {
			continue
		}

		firstNonDelimiter := true
		for _, v := range (*p.tokenizer.tokens)[:tokenIndex] {
			if v.Enclosed || v.Category == tokenCategoryDelimiter {
				continue
			} else {
				firstNonDelimiter = false
			}
		}
		if firstNonDelimiter {
			continue
		}

		previousToken, _ := p.tokenizer.tokens.findPrevious(*tkn, tokenFlagsNotDelimiter)
		if previousToken.Category == tokenCategoryUnknown {
			if strings.ToUpper(previousToken.Content) == "MOVIE" || strings.ToUpper(previousToken.Content) == "PART" {
				continue
			}
		}
		if p.setEpisodeNumber(tkn.Content, tkn, true) {
			return true
		}
	}
	return false
}

func (p *parser) setSeasonNumber(number string, tkn *token) bool {
	if !isNumeric(number) {
		return false
	}
	p.tokenizer.elements.insert(elementCategoryAnimeSeason, number)
	tkn.Category = tokenCategoryIdentifier
	return true
}

func isValidEpisodeNumber(number string) bool {
	return stringToInt(number) <= episodeNumberMax
}

func (p *parser) setEpisodeNumber(number string, tkn *token, validate bool) bool {
	if validate {
		if !isValidEpisodeNumber(number) {
			return false
		}
	}
	tkn.Category = tokenCategoryIdentifier
	cat := elementCategoryEpisodeNumber

	if p.tokenizer.elements.getCheckAltNumber() {
		episodeNumber := p.tokenizer.elements.get(elementCategoryEpisodeNumber)[0]
		if stringToInt(number) > stringToInt(episodeNumber) {
			cat = elementCategoryEpisodeNumberAlt
		} else if stringToInt(number) < stringToInt(episodeNumber) {
			p.tokenizer.elements.remove(elementCategoryEpisodeNumber, episodeNumber)
			p.tokenizer.elements.insert(elementCategoryEpisodeNumberAlt, episodeNumber)
		} else {
			return false
		}
	}

	p.tokenizer.elements.insert(cat, number)
	return true
}

func (p *parser) setAlternativeEpisodeNumber(number string, tkn *token) {
	p.tokenizer.elements.insert(elementCategoryEpisodeNumberAlt, number)
	tkn.Category = tokenCategoryIdentifier
}

func (p *parser) matchEpisodePattern(w string, tkn *token) bool {
	var numericFront bool
	var numericBack bool

	if isNumeric(w) {
		return false
	}

	w = strings.Trim(w, " -")

	if len(w) == 0 {
		return false
	}

	if isNumeric(string(w[0])) {
		numericFront = true
	}
	if isNumeric(string(w[len(w)-1])) {
		numericBack = true
	}

	if numericFront && numericBack {
		if p.matchSingleEpisodePattern(w, tkn) {
			return true
		}
	}
	if numericFront && numericBack {
		if p.matchMultiEpisodePattern(w, tkn) {
			return true
		}
	}
	if numericBack {
		if p.matchSeasonAndEpisodePattern(w, tkn) {
			return true
		}
	}
	if !numericFront {
		if p.matchTypeAndEpisodePattern(w, tkn) {
			return true
		}
	}
	if numericFront && numericBack {
		if p.matchFractionalEpisodePattern(w, tkn) {
			return true
		}
	}
	if numericFront && !numericBack {
		if p.matchPartialEpisodePattern(w, tkn) {
			return true
		}
	}
	if numericBack {
		if p.matchNumberSignPattern(w, tkn) {
			return true
		}
	}
	if numericFront || strings.IndexRune(w, '\u7B2C') == 0 {
		if p.matchJapaneseCounterPattern(w, tkn) {
			return true
		}
	}

	return false
}

func (p *parser) matchSingleEpisodePattern(w string, tkn *token) bool {
	pattern := "(\\d{1,4})[vV](\\d)$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	p.setEpisodeNumber(match[1], tkn, false)
	_, err := strconv.Atoi(match[2])
	if err == nil {
		p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[2])
	}

	return true
}

func (p *parser) matchMultiEpisodePattern(w string, tkn *token) bool {
	pattern := "(\\d{1,4})(?:[vV](\\d))?[-~&+](\\d{1,4})(?:[vV](\\d))?$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	lowerBound, _ := strconv.Atoi(match[1])
	upperBound, _ := strconv.Atoi(match[3])
	if lowerBound < upperBound {
		if p.setEpisodeNumber(match[1], tkn, true) {
			p.setEpisodeNumber(match[3], tkn, false)
			if len(match[2]) > 0 {
				p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[2])
			}
			if len(match[4]) > 0 {
				p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[4])
			}
			return true
		}
	}
	return false
}

func (p *parser) matchSeasonAndEpisodePattern(w string, tkn *token) bool {
	pattern := "(?i)S?(\\d{1,2})(?:-S?(\\d{1,2}))?(?:x|[ ._-x]?E)(\\d{1,4})(?:-E?(\\d{1,4}))?(?:[vV](\\d))?$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}

	p.tokenizer.elements.insert(elementCategoryAnimeSeason, match[1])
	if len(match[2]) > 0 {
		p.tokenizer.elements.insert(elementCategoryAnimeSeason, match[2])
	}
	p.setEpisodeNumber(match[3], tkn, false)
	if len(match[4]) > 0 {
		p.setEpisodeNumber(match[4], tkn, false)
	}
	if len(match[5]) > 0 {
		p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[5])
	}
	return true
}

func (p *parser) matchTypeAndEpisodePattern(w string, tkn *token) bool {
	numberBegin := findNumberInString(w)
	if numberBegin == -1 {
		return false
	}
	prefix := w[:numberBegin]

	kd, found := p.tokenizer.keywordManager.find(p.tokenizer.keywordManager.normalize(prefix), elementCategoryAnimeType)
	if found {
		p.tokenizer.elements.insert(elementCategoryAnimeType, prefix)
		number := w[numberBegin:]
		if p.matchEpisodePattern(number, tkn) || p.setEpisodeNumber(number, tkn, true) {
			tokenIndex := p.tokenizer.tokens.getIndex(*tkn, 0)
			tkn.Content = number
			targetCategory := tokenCategoryIdentifier
			if !kd.options.identifiable {
				targetCategory = tokenCategoryUnknown
			}
			p.tokenizer.tokens.insert(tokenIndex, token{
				Category: targetCategory,
				Content:  prefix,
				Enclosed: tkn.Enclosed,
			})
		}
		return true
	}
	return false
}

func (p *parser) matchFractionalEpisodePattern(w string, tkn *token) bool {
	pattern := "\\d+\\.5$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	if p.setEpisodeNumber(w, tkn, true) {
		return true
	}
	return false
}

func (p *parser) matchPartialEpisodePattern(w string, tkn *token) bool {
	nonNumberBegin := findNonNumberInString(w)
	if nonNumberBegin == -1 {
		return false
	}
	suffix := string(w[nonNumberBegin:])

	if len(suffix) == 1 && strings.Contains("ABCabc", suffix) {
		if p.setEpisodeNumber(w, tkn, true) {
			return true
		}
	}
	return false
}

func (p *parser) matchNumberSignPattern(w string, tkn *token) bool {
	if string(w[0]) != "#" {
		return false
	}

	pattern := "#(\\d{1,4})(?:[-~&+](\\d{1,4}))?(?:[vV](\\d))?$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	p.setEpisodeNumber(match[1], tkn, false)
	if len(match[2]) > 0 {
		p.setEpisodeNumber(match[2], tkn, true)
	}
	if len(match[3]) > 0 {
		p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[3])
	}
	return true
}

func (p *parser) matchJapaneseCounterPattern(w string, tkn *token) bool {
	if strings.IndexRune(w, '\u8A71') == -1 {
		return false
	}

	pattern := "(\\d{1,4})話$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	p.setEpisodeNumber(match[1], tkn, false)
	return true
}

func (p *parser) matchVolumePattern(w string, tkn *token) bool {
	var numericFront bool
	var numericBack bool

	if isNumeric(w) {
		return false
	}

	w = strings.Trim(w, " -")

	if isNumeric(string(w[0])) {
		numericFront = true
	}
	if isNumeric(string(w[len(w)-1])) {
		numericBack = true
	}

	if numericFront && numericBack {
		if p.matchSingleVolumePattern(w, tkn) {
			return true
		}
	}
	if numericFront && numericBack {
		if p.matchMultiVolumePattern(w, tkn) {
			return true
		}
	}
	return false
}

func isValidVolumeNumber(number string) bool {
	i, err := strconv.Atoi(number)
	if err != nil {
		return false
	}
	return i <= volumeNumberMax
}

func (p *parser) setVolumeNumber(number string, tkn *token, validate bool) bool {
	if validate {
		if !isValidVolumeNumber(number) {
			return false
		}
	}
	p.tokenizer.elements.insert(elementCategoryVolumeNumber, number)
	tkn.Category = tokenCategoryIdentifier
	return true
}

func (p *parser) matchSingleVolumePattern(w string, tkn *token) bool {
	pattern := "(\\d{1,2})[vV](\\d)$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	p.setVolumeNumber(match[1], tkn, false)
	p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[2])

	return true
}

func (p *parser) matchMultiVolumePattern(w string, tkn *token) bool {
	pattern := "(\\d{1,2})[-~&+](\\d{1,2})(?:[vV](\\d))?$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	lowerBound, _ := strconv.Atoi(match[1])
	upperBound, _ := strconv.Atoi(match[2])
	if lowerBound < upperBound {
		if p.setVolumeNumber(match[1], tkn, true) {
			p.setVolumeNumber(match[2], tkn, false)
			if len(match[3]) > 0 {
				p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[3])
			}
			return true
		}
	}
	return false
}
