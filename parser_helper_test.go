package anitogo

import (
	"testing"
)

func TestParserHelperCheckAnimeSeasonKeyword(t *testing.T) {
	testTkn := token{
		Content:  "",
		Category: tokenCategoryUnknown,
	}
	psr := getTestParser("")
	ret := psr.checkAnimeSeasonKeyword(&testTkn)
	if ret {
		t.Error("expected false, got true")
	}
	var testTkn2 *token
	psr = getTestParser("[Conclave-Mendoi]_Mobile_Suit_Gundam_00_S2_-_01v2_[1280x720_H.264_AAC][4863FBE8].mkv")
	for _, v := range *psr.tokenizer.tokens {
		if v.Content == "Gundam" {
			testTkn2 = v
		}
	}
	ret = psr.checkAnimeSeasonKeyword(testTkn2)
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserHelperSetAnimeSeason(t *testing.T) {
	psr := getTestParser("[Conclave-Mendoi]_Mobile_Suit_Gundam_00_S2_-_01v2_[1280x720_H.264_AAC][4863FBE8].mkv")
	firstTkn := (*psr.tokenizer.tokens)[0]
	secondTkn := (*psr.tokenizer.tokens)[1]
	content := "S2"
	psr.setAnimeSeason(firstTkn, secondTkn, content)
	if firstTkn.Category != tokenCategoryIdentifier {
		t.Errorf("expected %d, got %d", tokenCategoryIdentifier, firstTkn.Category)
	}
	if secondTkn.Category != tokenCategoryIdentifier {
		t.Errorf("expected %d, got %d", tokenCategoryIdentifier, secondTkn.Category)
	}
	if !psr.tokenizer.elements.contains(elementCategoryAnimeSeason) {
		t.Error("expected elements to contain anime season")
	}
}

func TestParserHelperBuildElement(t *testing.T) {
	psr := getTestParser("")
	firstTkn := (*psr.tokenizer.tokens)[3]
	secondTkn := (*psr.tokenizer.tokens)[5]
	psr.buildElement(elementCategoryAnimeTitle, firstTkn, secondTkn, true)
	if psr.tokenizer.elements.get(elementCategoryAnimeTitle)[0] != "_Toradora!_" {
		t.Errorf("expected \"_Toradora!_\", got %s", psr.tokenizer.elements.get(elementCategoryAnimeTitle)[0])
	}
	psr = getTestParser("")
	firstTkn = (*psr.tokenizer.tokens)[3]
	secondTkn = (*psr.tokenizer.tokens)[5]
	psr.buildElement(elementCategoryAnimeTitle, firstTkn, secondTkn, false)
	if psr.tokenizer.elements.get(elementCategoryAnimeTitle)[0] != "Toradora!" {
		t.Errorf("expected \"Toradora!\", got %s", psr.tokenizer.elements.get(elementCategoryAnimeTitle)[0])
	}
}

func TestParserHelperFindNonNumberInString(t *testing.T) {
	i := findNonNumberInString("1111")
	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}
	i = findNonNumberInString("111a")
	if i != 3 {
		t.Errorf("expected 3, got %d", i)
	}
}

func TestParserHelperIsDashCharacter(t *testing.T) {
	ret := isDashCharacter("")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isDashCharacter("1")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isDashCharacter("-")
	if !ret {
		t.Error("exptected true, got false")
	}
}

func TestParserHelperIsMostlyLatinString(t *testing.T) {
	ret := isMostlyLatinString("")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isMostlyLatinString("tt話")
	if !ret {
		t.Error("expected true, got false")
	}
	ret = isMostlyLatinString("t話話")
	if ret {
		t.Error("expected false, got true")
	}
}

func TestParserHelperStringToInt(t *testing.T) {
	i := stringToInt("9.5")
	if i != 9 {
		t.Errorf("expected 9, got %d", i)
	}
	i = stringToInt("test")
	if i != 0 {
		t.Errorf("expected 0, got %d", i)
	}
	i = stringToInt("8")
	if i != 8 {
		t.Errorf("expected 8, got %d", i)
	}
}

func TestParserHelperIsCRC32(t *testing.T) {
	ret := isCRC32("zzzzzzzz")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isCRC32("1234ABC")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isCRC32("1234ABCD")
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserHelperIsHexadecimalString(t *testing.T) {
	ret := isHexadecimalString("ZZZZ")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isHexadecimalString("1234ABCD")
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserHelperIsResolution(t *testing.T) {
	ret := isResolution("testxtest")
	if ret {
		t.Error("expected false, got true")
	}
	ret = isResolution("1920x1080")
	if !ret {
		t.Error("expected true, got false")
	}
	ret = isResolution("1080p")
	if !ret {
		t.Error("expected true, got false")
	}
}

func TestParserHelperGetNumberFromOrdinal(t *testing.T) {
	i := getNumberFromOrdinal("FIRST")
	if i != 1 {
		t.Errorf("expected 1, got %d", i)
	}
	i = getNumberFromOrdinal("One gazillion")
	if i != 0 {
		t.Errorf("expected 0, got %d", i)
	}
}

func TestParserHelperFindNumberInString(t *testing.T) {
	i := findNumberInString("aaa")
	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}
	i = findNumberInString("a1a")
	if i != 1 {
		t.Errorf("expected 1, got %d", i)
	}
}
