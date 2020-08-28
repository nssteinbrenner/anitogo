package anitogo

import (
	"testing"
)

func TestKeywordAdd(t *testing.T) {
	kwm := &keywordManager{
		keywords:       make(map[string]keyword),
		fileExtensions: make(map[string]keyword),
	}
	targKeyword1 := keyword{
		category: elementCategoryEpisodePrefix,
		options:  keywordOptionsDefault,
	}
	targKeyword2 := keyword{
		category: elementCategoryFileExtension,
		options:  keywordOptionsDefault,
	}

	kwm.add(elementCategoryEpisodePrefix, keywordOptionsDefault, []string{"TEST"})
	kwm.add(elementCategoryFileExtension, keywordOptionsDefault, []string{"TEST"})
	if kwm.keywords["TEST"] != targKeyword1 {
		t.Error("keyword did not match targKeyword1")
	}
	if kwm.fileExtensions["TEST"] != targKeyword2 {
		t.Error("keyword did not match targKeyword2")
	}
}

func TestKeywordFind(t *testing.T) {
	kwm := newKeywordManager()
	_, found := kwm.find("MKV", elementCategoryFileExtension)
	if !found {
		t.Error("expected true, got false")
	}
	_, found = kwm.find("MKV", elementCategoryUnknown)
	if found {
		t.Error("expected false, got true")
	}
	_, found = kwm.find("SAISON", elementCategoryAnimeSeasonPrefix)
	if !found {
		t.Error("expected true, got false")
	}
	_, found = kwm.find("SAISON", elementCategoryFileExtension)
	if found {
		t.Error("expected false, got true")
	}
}

func TestKeywordFindWithoutCategory(t *testing.T) {
	kwm := newKeywordManager()
	_, found := kwm.findWithoutCategory("MKV")
	if !found {
		t.Error("expected true, got false")
	}
	_, found = kwm.findWithoutCategory("SAISON")
	if !found {
		t.Error("expected true, got false")
	}
	_, found = kwm.findWithoutCategory("TESTINGTHIS")
	if found {
		t.Error("expected false, got true")
	}
}

func TestKeywordPeek(t *testing.T) {
	psr := getTestParser("")
	testStr := "this is a Dual Audio"
	idxSets := psr.tokenizer.keywordManager.peek(testStr, psr.tokenizer.elements)
	if idxSets[0].beginPos != 10 {
		t.Errorf("expected 10, got %d", idxSets[0].beginPos)
	}
	if idxSets[0].endPos != 20 {
		t.Errorf("expected 19, got %d", idxSets[0].endPos)
	}
	if testStr[idxSets[0].beginPos:idxSets[0].endPos] != "Dual Audio" {
		t.Errorf("expected \"%s\", got \"%s\"", "Dual Audio", testStr[idxSets[0].beginPos:idxSets[0].endPos])
	}
}
