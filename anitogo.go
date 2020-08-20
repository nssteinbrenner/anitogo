package anitogo

import (
	"strings"
	"unicode"
)

const maxExtensionLength = 4

// DefaultOptions to be passed to the Parse function.
//
// Custom options can be specified by creating a new Options struct.
//
// All testing was done with the default options, and they are recommended for most use cases.
var DefaultOptions = Options{
	AllowedDelimiters:  " _.&+,|",
	IgnoredStrings:     []string{},
	ParseEpisodeNumber: true,
	ParseEpisodeTitle:  true,
	ParseFileExtension: true,
	ParseReleaseGroup:  true,
}

// Parse builds and returns an Elements struct and an error by parsing a filename with the specified options.
//
// If an error is encountered during the process, Parse will return an empty Elements struct.
//
// Parsing behavior can be changed in the passed Options struct, which will change the returned Elements struct.
func Parse(filename string, options Options) (*Elements, error) {
	tkns := tokens{}
	elems := newElements()
	km := newKeywordManager()

	if len(filename) == 0 {
		return &Elements{}, nil
	}

	elems.insert(elementCategoryFileName, filename)
	newFilename, extension := removeExtensionFromFilename(km, filename)
	if newFilename != "" {
		filename = newFilename
	}
	if extension != "" {
		elems.insert(elementCategoryFileExtension, extension)
	}

	if options.IgnoredStrings != nil {
		filename = removeIgnoredStrings(filename, options.IgnoredStrings)
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
		return &Elements{}, err
	}

	psr := newParser(&tkz)
	err = psr.parse()
	if err != nil {
		return &Elements{}, err
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

func removeIgnoredStrings(filename string, ignoredStrings []string) string {
	for _, s := range ignoredStrings {
		filename = strings.ReplaceAll(filename, s, "")
	}
	return filename
}

func isAlphaNumeric(s string) bool {
	for _, v := range s {
		if !unicode.IsLetter(v) && !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
