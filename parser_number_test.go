package anitogo

import (
	"testing"
)

func TestParserNumberCheckExtentKeyword(t *testing.T) {
	psr := getTestParser("")
	for _, v := range *psr.tokenizer.tokens {
		v.Category = tokenCategoryIdentifier
	}
	ret := psr.checkExtentKeyword(elementCategoryEpisodeNumber, (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	(*psr.tokenizer.tokens)[0].Category = tokenCategoryUnknown
	(*psr.tokenizer.tokens)[1].Category = tokenCategoryUnknown
	(*psr.tokenizer.tokens)[1].Content = "11111"
	ret = psr.checkExtentKeyword(elementCategoryAnimeTitle, (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	var testTkn *token
	psr = getTestParser("")
	for _, v := range *psr.tokenizer.tokens {
		if v.Content == "-" {
			testTkn = v
			break
		}
	}
	ret = psr.checkExtentKeyword(elementCategoryEpisodeNumber, testTkn)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberSearchForEpisodePatterns(t *testing.T) {
	psr := getTestParser("")
	ret := psr.searchForEpisodePatterns(*psr.tokenizer.tokens)
	if !ret {
		t.Error("expected true, got false")
	}
	psr = getTestParser("[TaigaSubs]_Toradora!_(2008)_-_Tiger_and_Dragon_[1280x720_H.264_FLAC][1234ABCD].mkv")
	ret = psr.searchForEpisodePatterns(*psr.tokenizer.tokens)
	if ret {
		t.Error("expected false, got true")
	}
}

func TestParserNumberNumberComesAfterPrefix(t *testing.T) {
	psr := getTestParser("")
	var testTkn *token
	for _, v := range *psr.tokenizer.tokens {
		if v.Content == "01v2" {
			testTkn = v
		}
	}
	ret := psr.numberComesAfterPrefix(elementCategoryEpisodePrefix, testTkn)
	if ret {
		t.Error("expected false, got true")
	}
	psr = getTestParser("[Elysium]Sora.no.Woto.EP07.5(BD.720p.AAC)[C37580F8].mkv")
	for _, v := range *psr.tokenizer.tokens {
		if v.Content == "EP07.5" {
			testTkn = v
			break
		}
	}
	ret = psr.numberComesAfterPrefix(elementCategoryEpisodePrefix, testTkn)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberNumberComesBeforeAnotherNumber(t *testing.T) {
	psr := getTestParser("")
	var testTkn *token
	for _, v := range *psr.tokenizer.tokens {
		if v.Content == "2008" {
			testTkn = v
		}
	}

	ret := psr.numberComesBeforeAnotherNumber(testTkn)
	if ret {
		t.Error("expected false, got true")
	}
	psr = getTestParser("[TaigaSubs]_Toradora!_(2008_&_2009)_-_Tiger_and_Dragon_[1280x720_H.264_FLAC][1234ABCD].mkv")
	for _, v := range *psr.tokenizer.tokens {
		if v.Content == "2008" {
			testTkn = v
		}
	}
	ret = psr.numberComesBeforeAnotherNumber(testTkn)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberSearchForEquivalentNumbers(t *testing.T) {
	psr := getTestParser("")
	ret := psr.searchForEquivalentNumbers(*psr.tokenizer.tokens)
	if ret {
		t.Error("expected false, got true")
	}
	psr = getTestParser("__BLUE DROP 10 (1).avi")
	ret = psr.searchForEquivalentNumbers(*psr.tokenizer.tokens)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberSearchForSeparatedNumbers(t *testing.T) {
	psr := getTestParser("")
	ret := psr.searchForSeparatedNumbers(*psr.tokenizer.tokens)
	if ret {
		t.Error("expected false, got true")
	}
	psr = getTestParser("[ANBU]_Princess_Lover!_-_01_[2048A39A]")
	var testIdx int
	for i, v := range *psr.tokenizer.tokens {
		if v.Content == "01" {
			testIdx = i
			break
		}
	}
	testTkns := (*psr.tokenizer.tokens)[testIdx:]
	ret = psr.searchForSeparatedNumbers(testTkns)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberSearchForIsolatedNumbersTokens(t *testing.T) {
	psr := getTestParser("1")
	ret := psr.searchForIsolatedNumbersTokens(*psr.tokenizer.tokens)
	if ret {
		t.Error("expected false, got true")
	}
	psr = getTestParser("")
	ret = psr.searchForIsolatedNumbersTokens(*psr.tokenizer.tokens)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberSearchForLastNumber(t *testing.T) {
	psr := getTestParser("")
	for _, v := range *psr.tokenizer.tokens {
		v.Enclosed = true
	}
	ret := psr.searchForLastNumber(*psr.tokenizer.tokens)
	if ret {
		t.Error("expected false, got true")
	}
	psr = getTestParser("")
	ret = psr.searchForLastNumber(*psr.tokenizer.tokens)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberSetSeasonNumber(t *testing.T) {
	psr := getTestParser("")
	ret := psr.setSeasonNumber("test", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	if psr.tokenizer.elements.contains(elementCategoryAnimeSeason) {
		t.Error("expected false, got true")
	}

	psr = getTestParser("")
	ret = psr.setSeasonNumber("1", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
	if !psr.tokenizer.elements.contains(elementCategoryAnimeSeason) {
		t.Error("expected true, got false")
	}
	if (*psr.tokenizer.tokens)[0].Category != tokenCategoryIdentifier {
		t.Errorf("expected %d, got %d", tokenCategoryIdentifier, (*psr.tokenizer.tokens)[0].Category)
	}
}

func TestParserNumberSetEpisodeNumber(t *testing.T) {
	psr := getTestParser("")
	psr.tokenizer.elements.insert(elementCategoryEpisodeNumber, "1")
	ret := psr.setEpisodeNumber("9999", (*psr.tokenizer.tokens)[0], true)
	if ret {
		t.Error("expected false, got true")
	}
	psr.tokenizer.elements.setCheckAltNumber(true)
	ret = psr.setEpisodeNumber("999", (*psr.tokenizer.tokens)[0], false)
	if !ret {
		t.Error("expected true, got false")
	}
	if (*psr.tokenizer.tokens)[0].Category != tokenCategoryIdentifier {
		t.Errorf("expected %d, got %d", tokenCategoryIdentifier, (*psr.tokenizer.tokens)[0].Category)
	}
}

func TestParserNumberSetAlternativeEpisodeNumber(t *testing.T) {
	psr := getTestParser("")
	psr.setAlternativeEpisodeNumber("1", (*psr.tokenizer.tokens)[0])
	if !psr.tokenizer.elements.contains(elementCategoryEpisodeNumberAlt) {
		t.Error("expected true, got false")
	}
	if (*psr.tokenizer.tokens)[0].Category != tokenCategoryIdentifier {
		t.Errorf("expected %d, got %d", tokenCategoryIdentifier, (*psr.tokenizer.tokens)[0].Category)
	}
}

func TestParserNumberMatchMultiEpisodePattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchMultiEpisodePattern("t01v2-05v2", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchMultiEpisodePattern("01v2-05v2", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchSeasonAndEpisodePattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchSeasonAndEpisodePattern("ts01e01", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchSeasonAndEpisodePattern("s01e01", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchFractionalEpisodePattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchFractionalEpisodePattern("t1.5", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	(*psr.tokenizer.tokens)[0].Content = "99999.5"
	ret = psr.matchFractionalEpisodePattern("99999.5", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchFractionalEpisodePattern("1.5", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchPartialEpisodePattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchPartialEpisodePattern("11111", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchPartialEpisodePattern("1A", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchNumberSignPattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchNumberSignPattern("#", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchNumberSignPattern("#t#1", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchNumberSignPattern("#1", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchJapaneseCounterPattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchJapaneseCounterPattern("話test", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchJapaneseCounterPattern("12話", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchVolumePattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchVolumePattern("test", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchVolumePattern("01v2", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberIsValidVolumeNumber(t *testing.T) {
	ret := isValidVolumeNumber("test")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isValidVolumeNumber("21")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isValidVolumeNumber("20")
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberSetVolumeNumber(t *testing.T) {
	psr := getTestParser("")
	ret := psr.setVolumeNumber("9999999", (*psr.tokenizer.tokens)[0], true)
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.setVolumeNumber("9", (*psr.tokenizer.tokens)[0], true)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchSingleVolumePattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchSingleVolumePattern("test1v2", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchSingleVolumePattern("1v2", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserNumberMatchMultiVolumePattern(t *testing.T) {
	psr := getTestParser("")
	ret := psr.matchMultiVolumePattern("test", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchMultiVolumePattern("test1-2v2", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchMultiVolumePattern("10-1v2", (*psr.tokenizer.tokens)[0])
	if ret {
		t.Error("expected false, got true")
	}
	ret = psr.matchMultiVolumePattern("1-2v2", (*psr.tokenizer.tokens)[0])
	if !ret {
		t.Error("expected true, got false")
	}
}

func getTestParser(filename string) *parser {
	if filename == "" {
		filename = "[TaigaSubs]_Toradora!_(2008)_-_01v2_-_Tiger_and_Dragon_[1280x720_H.264_FLAC][1234ABCD].mkv"
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

	tkz := tokenizer{
		filename:       filename,
		options:        DefaultOptions,
		tokens:         tkns,
		keywordManager: km,
		elements:       elems,
	}
	tkz.tokenize()

	psr := newParser(&tkz)

	return psr
}
