package anitogo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Options is a struct that allows you to change the parsing behavior.
//
// Default options have been provided under a variable named "DefaultOptions".
type Options struct {
	// DefaultOptions value: " _.&+,|"
	// Each character in this string will be evaluated as a delimiter during parsing.
	// The defaults are fairly sane, but in some cases you may want to change them.
	// For example in the following filename: DRAMAtical Murder Episode 1 - Data_01_Login
	// With the defaults, the "_" characters would be replaced with spaces, but this may
	// not be desired behavior.
	AllowedDelimiters string

	// DefaultOptions value: []string{}
	// These strings will be removed from the filename.
	IgnoredStrings []string

	// DefaultOptions value: true
	// Determines if the episode number will be parsed into the Elements struct.
	ParseEpisodeNumber bool

	// DefaultOptions value: true
	// Determines if the episode title will be parsed into the Elements struct.
	ParseEpisodeTitle bool

	// DefaultOptions value: true
	// Determines if the file extension will be parsed into the Elements struct.
	ParseFileExtension bool

	// DefaultOptions value: true
	// Determines if the release group will be parsed into the Elements struct.
	ParseReleaseGroup bool
}

type tokenizer struct {
	filename       string
	options        Options
	tokens         *tokens
	keywordManager *keywordManager
	elements       *Elements
}

func (t *tokenizer) addToken(cat int, content string, enclosed bool) {
	t.tokens.appendToken(token{
		Category: cat,
		Content:  content,
		Enclosed: enclosed,
	})
}

func (t *tokenizer) tokenize() {
	brackets := [][]rune{
		{'(', ')'},
		{'[', ']'},
		{'{', '}'},
		{'\u300C', '\u300D'},
		{'\u300E', '\u300F'},
		{'\u3010', '\u3011'},
		{'\uFF08', '\uFF09'},
	}

	text := t.filename
	isBracketOpen := false
	var matchingBracket rune
	for len(text) > 0 {
		var bracketIndex int
		if !isBracketOpen {
			bracketIndex, matchingBracket = findFirstBracket(text, brackets)
		} else {
			bracketIndex = strings.IndexRune(text, matchingBracket)
		}

		if bracketIndex != 0 {
			if bracketIndex != -1 {
				t.tokenizeByPreidentified(text[:bracketIndex], isBracketOpen)
			} else {
				t.tokenizeByPreidentified(text, isBracketOpen)
			}
		}

		if bracketIndex != -1 {
			t.addToken(tokenCategoryBracket, string(text[bracketIndex]), true)
			isBracketOpen = !isBracketOpen
			text = text[bracketIndex+1:]
		} else {
			text = ""
		}
	}
}

func (t *tokenizer) tokenizeByPreidentified(filename string, enclosed bool) {
	preIdentifiedtokens := t.keywordManager.peek(filename, t.elements)

	lastTokenEndPos := 0
	for _, preIdentified := range preIdentifiedtokens {
		tknBeginPos := preIdentified.beginPos
		tknEndPos := preIdentified.endPos
		if lastTokenEndPos != tknBeginPos && tknBeginPos <= len(filename) {
			t.tokenizeByDelimiters(filename[lastTokenEndPos:tknBeginPos], enclosed)
		}
		if tknEndPos <= len(filename) {
			t.addToken(tokenCategoryIdentifier, filename[tknBeginPos:tknEndPos], enclosed)
			lastTokenEndPos = tknEndPos
		}
	}
	if lastTokenEndPos != len(filename) {
		t.tokenizeByDelimiters(filename[lastTokenEndPos:], enclosed)
	}
}

func (t *tokenizer) tokenizeByDelimiters(filename string, enclosed bool) {
	var delimiters string
	var splitText []string
	for _, delimiter := range t.options.AllowedDelimiters {
		delimiters = delimiters + "\\" + string(delimiter)
	}
	pattern := fmt.Sprintf("([%v])", delimiters)
	text := filename
	re := regexp.MustCompile(pattern)
	splitText = splitWith(re, text, -1)
	for _, subtext := range splitText {
		if subtext != "" {
			if strings.Contains(t.options.AllowedDelimiters, subtext) {
				t.addToken(tokenCategoryDelimiter, subtext, enclosed)
			} else {
				t.addToken(tokenCategoryUnknown, subtext, enclosed)
			}
		}
	}
	t.validateDelimitertokens()
}

func (t *tokenizer) validateDelimitertokens() {
	for _, tkn := range *t.tokens {
		if tkn.Category != tokenCategoryDelimiter {
			continue
		}
		delimiter := tkn.Content
		prevToken, _ := t.findPreviousValidToken(tkn)
		nextToken, _ := t.findNextValidToken(tkn)

		if delimiter != " " && delimiter != "_" {
			if t.isSingleCharacterToken((*prevToken)) {
				nestedNextToken := *nextToken
				prevToken = t.appendTokenTo(tkn, prevToken)
				for t.isUnknownToken(nestedNextToken) {
					prevToken = t.appendTokenTo(&nestedNextToken, prevToken)
					if nestedNextToken.Content == nextToken.Content {
						nextToken.Category = tokenCategoryInvalid
					}
					holder, _ := t.findNextValidToken(&nestedNextToken)
					nestedNextToken = *holder
					if t.isDelimiterToken(nestedNextToken) && nestedNextToken.Content == delimiter {
						prevToken = t.appendTokenTo(&nestedNextToken, prevToken)
						holder, _ = t.findNextValidToken(&nestedNextToken)
						nestedNextToken = *holder
					}
					continue
				}
			}
			if t.isSingleCharacterToken((*nextToken)) {
				prevToken = t.appendTokenTo(tkn, prevToken)
				t.appendTokenTo(nextToken, prevToken)
				continue
			}
		}

		if t.isUnknownToken((*prevToken)) && t.isDelimiterToken((*nextToken)) {
			nextDelimiter := nextToken.Content
			if delimiter != nextDelimiter && delimiter != "," {
				if nextDelimiter == " " || nextDelimiter == "_" {
					prevToken = t.appendTokenTo(tkn, prevToken)
				}
			}
		} else if t.isDelimiterToken((*prevToken)) && t.isDelimiterToken((*nextToken)) {
			prevDelimiter := prevToken.Content
			nextDelimiter := nextToken.Content
			if prevDelimiter == nextDelimiter && prevDelimiter != delimiter {
				tkn.Category = tokenCategoryUnknown
			}
		}
		if delimiter == "&" || delimiter == "+" {
			if t.isUnknownToken((*prevToken)) && t.isUnknownToken((*nextToken)) {
				if isNumeric(prevToken.Content) && isNumeric(nextToken.Content) {
					prevToken = t.appendTokenTo(tkn, prevToken)
					t.appendTokenTo(nextToken, prevToken)
				}
			}
		}
	}
	var newTkns tokens
	for _, tkn := range *t.tokens {
		if tkn.Category != tokenCategoryInvalid {
			newTkns = append(newTkns, tkn)
		}
	}
	t.tokens.update(newTkns)
}

func (t *tokenizer) findPreviousValidToken(tkn *token) (*token, bool) {
	return t.tokens.findPrevious(*tkn, tokenFlagsValid)
}

func (t *tokenizer) findNextValidToken(tkn *token) (*token, bool) {
	return t.tokens.findNext(*tkn, tokenFlagsValid)
}

func (t *tokenizer) isDelimiterToken(tkn token) bool {
	if !tkn.empty() && tkn.Category == tokenCategoryDelimiter {
		return true
	}
	return false
}

func (t *tokenizer) isUnknownToken(tkn token) bool {
	if !tkn.empty() && tkn.Category == tokenCategoryUnknown {
		return true
	}
	return false
}

func (t *tokenizer) isSingleCharacterToken(tkn token) bool {
	if t.isUnknownToken(tkn) && len(tkn.Content) == 1 && tkn.Content != "-" {
		return true
	}
	return false
}

func (t *tokenizer) appendTokenTo(tkn, appendTo *token) *token {
	appendToIndex := t.tokens.getIndex(*appendTo, 0)
	appendToSrc, _ := t.tokens.get(appendToIndex)
	appendToSrc.Content += tkn.Content
	srcTknIndex := t.tokens.getIndex(*tkn, appendToIndex)
	srcTkn, _ := t.tokens.get(srcTknIndex)
	srcTkn.Category = tokenCategoryInvalid

	return appendToSrc
}

func findFirstBracket(filename string, brackets [][]rune) (int, rune) {
	var openBrackets []rune
	for _, v := range brackets {
		openBrackets = append(openBrackets, v[0])
	}

	index := -1
	for idx, bracket := range filename {
		var found bool
		for _, v := range openBrackets {
			if bracket == v {
				found = true
				break
			}
		}
		if found {
			index = idx
			break
		}
	}

	var matchingBracket rune
	for _, v := range brackets {
		if index != -1 {
			if strings.IndexRune(filename, v[0]) == index {
				matchingBracket = v[1]
			}
		}
	}
	return index, matchingBracket
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func splitWith(re *regexp.Regexp, s string, n int) []string {
	if n == 0 {
		return nil
	}

	matches := re.FindAllStringIndex(s, n)
	strings := make([]string, 0, len(matches))

	beg := 0
	end := 0
	for _, match := range matches {
		if n > 0 && len(strings) >= n-1 {
			break
		}

		end = match[0]
		if match[1] != 0 {
			strings = append(strings, s[beg:end])
		}
		beg = match[1]
		strings = append(strings, s[match[0]:match[1]])
	}

	if end != len(s) {
		strings = append(strings, s[beg:])
	}

	return strings
}
