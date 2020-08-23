package anitogo

import (
	"strings"
	"unicode"
)

const maxExtensionLength = 4

// DefaultOptions are sane defaults for the Options struct to be passed to the Parse function.
//
// Custom options can be specified by creating a new Options struct and passing it to the Parse function.
var DefaultOptions = Options{
	AllowedDelimiters:  " _.&+,|",
	IgnoredStrings:     []string{},
	ParseEpisodeNumber: true,
	ParseEpisodeTitle:  true,
	ParseFileExtension: true,
	ParseReleaseGroup:  true,
}

// Parse builds and returns a pointer to an Elements struct and an error by parsing a filename with the specified options.
//
// If an error is encountered during the process, Parse will return an empty pointer to an Elements struct.
//
// Parsing behavior can be customized in the passed Options struct.
func Parse(filename string, options Options) (*Elements, error) {
	if len(filename) == 0 {
		return &Elements{}, traceError(emptyFilenameErr)
	}

	tkns := &tokens{}
	elems := &Elements{}
	km := newKeywordManager()

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
		tokens:         tkns,
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
	_, found := km.find(km.normalize(extension), elementCategoryFileExtension)
	if !found {
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
