package anitogo

import (
	"strings"
	"unicode"
)

const maxExtensionLength = 4

var DefaultOptions = Options{
	AllowedDelimiters:  " _.&+,|",
	IgnoredStrings:     []string{},
	ParseEpisodeNumber: true,
	ParseEpisodeTitle:  true,
	ParseFileExtension: true,
	ParseReleaseGroup:  true,
}

func Parse(filename string, options Options) (*elements, error) {
	tkns := tokens{}
	elems := newElements()
	km := newKeywordManager()

	if len(filename) == 0 {
		return &elements{}, nil
	}

	elems.insert(elementCategoryFileName, filename)
	newFilename, extension := removeExtensionFromFilename(km, filename)
	if newFilename != "" {
		filename = newFilename
	}
	if extension != "" {
		elems.insert(elementCategoryFileExtension, extension)
	}

	tkz := tokenizer{
		filename:       filename,
		options:        options,
		tokens:         &tkns,
		keywordManager: km,
		elements:       elems,
	}
	err := tkz.tokenize()
	if err != nil {
		return &elements{}, err
	}

	psr := newParser(&tkz)
	err = psr.parse()
	if err != nil {
		return &elements{}, err
	}

	return psr.tokenizer.elements, nil
}

func removeExtensionFromFilename(km *keywordManager, filename string) (string, string) {
	var extension string

	extStart := strings.LastIndex(filename, ".")
	if extStart == -1 {
		return "", ""
	}

	extension = filename[extStart+1:]
	if len(extension) > maxExtensionLength {
		return "", ""
	}
	if !isAlphaNumeric(extension) {
		return "", ""
	}
	kd, found := km.find(km.Normalize(extension), elementCategoryFileExtension)
	if !found {
		return "", ""
	}
	if kd.Category != elementCategoryFileExtension {
		return "", ""
	}

	filename = filename[:extStart]

	return filename, extension
}

func isAlphaNumeric(s string) bool {
	for _, v := range s {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
