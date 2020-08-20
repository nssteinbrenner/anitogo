package anitogo

import (
	"fmt"
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

func (p *parser) checkExtentKeyword(cat elementCategory, tkn *token) (bool, error) {
	nextToken, _, err := p.tokenizer.tokens.findNext(*tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, err
	}

	if nextToken.Category == tokenCategoryUnknown {
		if !nextToken.empty() && findNumberInString(nextToken.Content) > -1 {
			if cat == elementCategoryEpisodeNumber {
				match, err := p.matchEpisodePattern(nextToken.Content, nextToken)
				if err != nil {
					return false, err
				}
				if !match {
					p.setEpisodeNumber(nextToken.Content, nextToken, false)
				}
			} else if cat == elementCategoryVolumeNumber {
				if !p.matchVolumePattern(nextToken.Content, nextToken) {
					p.setVolumeNumber(nextToken.Content, nextToken, false)
				}
			} else {
				return false, nil
			}
			tkn.Category = tokenCategoryIdentifier
			return true, nil
		}
	}

	return false, nil
}

func (p *parser) searchForEpisodePatterns(tkns tokens) (bool, error) {
	for _, tkn := range tkns {
		numericFront := isNumeric(string(tkn.Content[0]))

		if !numericFront {
			match, err := p.numberComesAfterPrefix(elementCategoryEpisodePrefix, tkn)
			if err != nil {
				return false, err
			}
			if match {
				return true, nil
			}
			match, err = p.numberComesAfterPrefix(elementCategoryVolumePrefix, tkn)
			if err != nil {
				return false, err
			}
			if match {
				continue
			}
			match, err = p.numberComesAfterPrefix(elementCategoryAnimeSeasonPrefix, tkn)
			if err != nil {
				return false, err
			}
			if match {
				continue
			}
		} else {
			match, err := p.numberComesBeforeAnotherNumber(tkn)
			if err != nil {
				return false, err
			}
			if match {
				return true, nil
			}
		}
		match, err := p.matchEpisodePattern(tkn.Content, tkn)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (p *parser) numberComesAfterPrefix(cat elementCategory, tkn *token) (bool, error) {
	numberBegin := findNumberInString(tkn.Content)
	if numberBegin == -1 {
		return false, nil
	}
	prefix := tkn.Content[:numberBegin]

	_, found := p.tokenizer.keywordManager.find(p.tokenizer.keywordManager.Normalize(prefix), cat)
	if found {
		number := tkn.Content[numberBegin:]
		if cat == elementCategoryEpisodePrefix {
			match, err := p.matchEpisodePattern(number, tkn)
			if err != nil {
				return false, err
			}
			if match {
				return true, nil
			}
			return p.setEpisodeNumber(number, tkn, false), nil
		}
		if cat == elementCategoryVolumePrefix {
			if p.matchVolumePattern(number, tkn) {
				return true, nil
			}
			return p.setVolumeNumber(number, tkn, false), nil
		}
		if cat == elementCategoryAnimeSeasonPrefix {
			return p.setSeasonNumber(number, tkn), nil
		}
	}
	return false, nil
}

func (p *parser) numberComesBeforeAnotherNumber(tkn *token) (bool, error) {
	separatorToken, found, err := p.tokenizer.tokens.findNext(*tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, err
	}

	if found {
		separator := separatorToken.Content
		if separator == "&" || separator == "of" {
			otherToken, found, err := p.tokenizer.tokens.findNext(*separatorToken, tokenFlagsNotDelimiter)
			if err != nil {
				return false, err
			}
			if found && isNumeric(otherToken.Content) {
				p.setEpisodeNumber(tkn.Content, tkn, false)
				if separator == "&" {
					p.setEpisodeNumber(otherToken.Content, tkn, false)
				}
				separatorToken.Category = tokenCategoryIdentifier
				otherToken.Category = tokenCategoryIdentifier
				return true, nil
			}
		}
	}
	return false, nil
}

func (p *parser) searchForEquivalentNumbers(tkns tokens) (bool, error) {
	for _, tkn := range tkns {
		iso, err := p.tokenizer.tokens.isTokenIsolated(*tkn)
		if err != nil {
			return false, err
		}
		if iso || !isValidEpisodeNumber(tkn.Content) {
			return false, err
		}

		nextToken, found, err := p.tokenizer.tokens.findNext(*tkn, tokenFlagsNotDelimiter)
		if err != nil {
			return false, err
		}
		if nextToken.empty() || nextToken.Category != tokenCategoryBracket || !found {
			continue
		}
		nextToken, found, err = p.tokenizer.tokens.findNext(*nextToken, tokenFlagsEnclosed|tokenFlagsNotDelimiter)
		if err != nil {
			return false, err
		}
		if found {
			if nextToken.Category != tokenCategoryUnknown {
				continue
			}
		}

		iso, err = p.tokenizer.tokens.isTokenIsolated(*nextToken)
		if err != nil {
			return false, err
		}
		if !iso || !isNumeric(nextToken.Content) || !isValidEpisodeNumber(nextToken.Content) {
			continue
		}

		i, err := strconv.Atoi(nextToken.Content)
		if err != nil {
			continue
		}
		j, err := strconv.Atoi(tkn.Content)
		if err != nil {
			continue
		}

		episode := nextToken
		altEpisode := tkn
		if i > j {
			episode = tkn
			altEpisode = nextToken
		}
		p.setEpisodeNumber(episode.Content, episode, false)
		p.setAlternativeEpisodeNumber(altEpisode.Content, altEpisode)
		return true, nil
	}

	return false, nil
}

func (p *parser) searchForSeparatedNumbers(tkns tokens) (bool, error) {
	for _, tkn := range tkns {
		previousToken, found, err := p.tokenizer.tokens.findPrevious(*tkn, tokenFlagsNotDelimiter)
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}

		if previousToken.Category == tokenCategoryUnknown && isDashCharacter(previousToken.Content) {
			if p.setEpisodeNumber(tkn.Content, tkn, true) {
				previousToken.Category = tokenCategoryIdentifier
				return true, nil
			}
		}
	}
	return false, nil
}

func (p *parser) searchForIsolatedNumbersTokens(tkns tokens) (bool, error) {
	for _, tkn := range tkns {
		iso, err := p.tokenizer.tokens.isTokenIsolated(*tkn)
		if err != nil {
			return false, err
		}
		if !tkn.Enclosed || !iso {
			continue
		}
		if p.setEpisodeNumber(tkn.Content, tkn, true) {
			return true, nil
		}
	}
	return false, nil
}

func (p *parser) searchForLastNumber(tkns tokens) (bool, error) {
	for _, tkn := range tkns {
		tokenIndex, err := p.tokenizer.tokens.getIndex(*tkn, 0)
		if err != nil {
			return false, err
		}

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

		previousToken, _, err := p.tokenizer.tokens.findPrevious(*tkn, tokenFlagsNotDelimiter)
		if err != nil {
			return false, err
		}
		if previousToken.Category == tokenCategoryUnknown {
			if strings.ToUpper(previousToken.Content) == "MOVIE" || strings.ToUpper(previousToken.Content) == "PART" {
				continue
			}
		}
		if p.setEpisodeNumber(tkn.Content, tkn, true) {
			return true, nil
		}
	}
	return false, nil
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
		episodeNumber := string(p.tokenizer.elements.get(elementCategoryEpisodeNumber)[0])
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

func (p *parser) setAlternativeEpisodeNumber(number string, tkn *token) bool {
	p.tokenizer.elements.insert(elementCategoryEpisodeNumberAlt, number)
	tkn.Category = tokenCategoryIdentifier

	return true
}

func (p *parser) matchEpisodePattern(w string, tkn *token) (bool, error) {
	var numericFront bool
	var numericBack bool

	if isNumeric(w) {
		return false, nil
	}

	w = strings.Trim(w, " -")

	if len(w) == 0 {
		return false, nil
	}

	if isNumeric(string(w[0])) {
		numericFront = true
	}
	if isNumeric(string(w[len(w)-1])) {
		numericBack = true
	}

	if numericFront && numericBack {
		if p.matchSingleEpisodePattern(w, tkn) {
			return true, nil
		}
	}
	if numericFront && numericBack {
		if p.matchMultiEpisodePattern(w, tkn) {
			return true, nil
		}
	}
	if numericBack {
		if p.matchSeasonAndEpisodePattern(w, tkn) {
			return true, nil
		}
	}
	if !numericFront {
		match, err := p.matchTypeAndEpisodePattern(w, tkn)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	if numericFront && numericBack {
		if p.matchFractionalEpisodePattern(w, tkn) {
			return true, nil
		}
	}
	if numericFront && !numericBack {
		if p.matchPartialEpisodePattern(w, tkn) {
			return true, nil
		}
	}
	if numericBack {
		if p.matchNumberSignPattern(w, tkn) {
			return true, nil
		}
	}
	if numericFront {
		if p.matchJapaneseCounterPattern(w, tkn) {
			return true, nil
		}
	}

	return false, nil
}

func (p *parser) matchSingleEpisodePattern(w string, tkn *token) bool {
	pattern := "(\\d{1,3})[vV](\\d)$"
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
	pattern := "(\\d{1,3})(?:[vV](\\d))?[-~&+](\\d{1,3})(?:[vV](\\d))?$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	lowerBound, err := strconv.Atoi(match[1])
	if err != nil {
		return false
	}
	upperBound, err := strconv.Atoi(match[3])
	if err != nil {
		return false
	}
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
	pattern := "(?i)S?(\\d{1,2})(?:-S?(\\d{1,2}))?(?:x|[ ._-x]?E)(\\d{1,3})(?:-E?(\\d{1,3}))?$"
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
	return true
}

func (p *parser) matchTypeAndEpisodePattern(w string, tkn *token) (bool, error) {
	numberBegin := findNumberInString(w)
	if numberBegin == -1 {
		return false, nil
	}
	prefix := w[:numberBegin]

	kd, found := p.tokenizer.keywordManager.find(p.tokenizer.keywordManager.Normalize(prefix), elementCategoryAnimeType)
	if found {
		p.tokenizer.elements.insert(elementCategoryAnimeType, prefix)
		number := w[numberBegin:]
		match, err := p.matchEpisodePattern(number, tkn)
		set := p.setEpisodeNumber(number, tkn, true)
		if err != nil {
			return false, err
		}
		if match || set {
			tokenIndex, err := p.tokenizer.tokens.getIndex(*tkn, 0)
			if err != nil {
				return false, err
			}
			if tokenIndex == -1 {
				return false, nil
			}
			tkn.Content = number
			targetCategory := tokenCategoryIdentifier
			if !kd.Options.Identifiable {
				targetCategory = tokenCategoryUnknown
			}
			err = p.tokenizer.tokens.insert(tokenIndex, token{
				Category: targetCategory,
				Content:  prefix,
				Enclosed: tkn.Enclosed,
			})
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, nil
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

	pattern := "#(\\d{1,3})(?:[-~&+](\\d{1,3}))?(?:[vV](\\d))?$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	if p.setEpisodeNumber(match[1], tkn, true) {
		if len(match[2]) > 0 {
			p.setEpisodeNumber(match[2], tkn, true)
		}
		if len(match[3]) > 0 {
			p.tokenizer.elements.insert(elementCategoryReleaseVersion, match[3])
		}
		return true
	}
	return false
}

func (p *parser) matchJapaneseCounterPattern(w string, tkn *token) bool {
	if string(w[len(w)-1]) != strings.Trim(fmt.Sprintf("%q", '\u8A71'), "'") {
		return false
	}

	pattern := "(\\d{1,3})\u8A71$"
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(w)
	if match == nil {
		return false
	}
	if strings.Index(w, match[0]) != 0 {
		return false
	}
	if p.setEpisodeNumber(match[1], tkn, false) {
		return true
	}

	return false
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
	lowerBound, err := strconv.Atoi(match[1])
	if err != nil {
		return false
	}
	upperBound, err := strconv.Atoi(match[2])
	if err != nil {
		return false
	}
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
