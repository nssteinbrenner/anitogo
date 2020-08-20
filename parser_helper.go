package anitogo

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const dashes = "-\u2010\u2011\u2012\u2013\u2014\u2015"

func stringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func isCRC32(str string) bool {
	return len(str) == 8 && isHexadecimalString(str)
}

func isHexadecimalString(str string) bool {
	_, err := strconv.ParseInt(str, 16, 64)
	return err == nil
}

func isResolution(str string) bool {
	pattern := "\\d{3,4}([pP]|([xX\u00D7]\\d{3,4}))$"
	found, _ := regexp.Match(pattern, []byte(str))
	return found
}

func getNumberFromOrdinal(str string) int {
	ordinals := map[string]int{
		"1st": 1, "first": 1,
		"2nd": 2, "second": 2,
		"3rd": 3, "third": 3,
		"4th": 4, "fourth": 4,
		"5th": 5, "fifth": 5,
		"6th": 6, "sixth": 6,
		"7th": 7, "seventh": 7,
		"8th": 8, "eigth": 8,
		"9th": 9, "ninth": 9,
	}

	lowerStr := strings.ToLower(str)
	num, _ := ordinals[lowerStr]
	return num
}

func findNumberInString(str string) int {
	for _, c := range str {
		if unicode.IsDigit(c) {
			return strings.IndexRune(str, c)
		}
	}
	return -1
}

func (p *parser) checkAnimeSeasonKeyword(tkn *token) (bool, error) {
	prevToken, found, err := p.tokenizer.tokens.findPrevious(*tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, err
	}
	if found {
		num := getNumberFromOrdinal(prevToken.Content)
		if num != 0 {
			p.setAnimeSeason(prevToken, tkn, strconv.Itoa(num))
			return true, nil
		}
	}

	nextToken, found, err := p.tokenizer.tokens.findNext(*tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, err
	}
	if found && isNumeric(nextToken.Content) {
		p.setAnimeSeason(tkn, nextToken, nextToken.Content)
		return true, nil
	}
	return false, nil
}

func (p *parser) setAnimeSeason(first, second *token, content string) error {
	p.tokenizer.elements.insert(elementCategoryAnimeSeason, content)
	firstIdx, err := p.tokenizer.tokens.getIndex(*first, 0)
	if err != nil {
		return err
	}
	secondIdx, err := p.tokenizer.tokens.getIndex(*second, firstIdx)
	if err != nil {
		return err
	}
	firstTkn, _ := p.tokenizer.tokens.get(firstIdx)
	secondTkn, _ := p.tokenizer.tokens.get(secondIdx)
	firstTkn.Category = tokenCategoryIdentifier
	secondTkn.Category = tokenCategoryIdentifier

	return nil
}

func (p *parser) buildElement(cat elementCategory, beginToken, endToken *token, keepDelimiters bool) error {
	element := ""

	tknList, err := p.tokenizer.tokens.getList(-1, beginToken, endToken)
	if err != nil {
		return err
	}
	for _, tkn := range tknList {
		if tkn.Category == tokenCategoryUnknown {
			element += tkn.Content
			tkn.Category = tokenCategoryIdentifier
		} else if tkn.Category == tokenCategoryBracket {
			element += tkn.Content
		} else if tkn.Category == tokenCategoryDelimiter {
			delimiter := tkn.Content
			if keepDelimiters {
				element += delimiter
			} else if tkn != beginToken && tkn != endToken {
				if delimiter == "," || delimiter == "&" {
					element += delimiter
				} else {
					element += " "
				}
			}
		}
	}

	if !keepDelimiters {
		element = strings.Trim(element, " "+dashes)
	}

	if element != "" {
		p.tokenizer.elements.insert(cat, strings.Trim(strings.ToValidUTF8(element, ""), " "))
	}
	return nil
}

func findNonNumberInString(str string) int {
	for _, r := range str {
		if !unicode.IsDigit(r) {
			return strings.IndexRune(str, r)
		}
	}
	return -1
}

func isDashCharacter(str string) bool {
	if len(str) != 1 {
		return false
	}
	for _, dash := range dashes {
		if str == string(dash) {
			return true
		}
	}
	return false
}

func isLatinRune(r rune) bool {
	return unicode.In(r, unicode.Latin)
}

func isMostlyLatinString(str string) bool {
	if len(str) <= 0 {
		return false
	}
	latinLength := 0
	for _, r := range str {
		if isLatinRune(r) {
			latinLength++
		}
	}
	return (float32(latinLength) / float32(len(str))) >= 0.5
}
