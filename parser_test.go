package anitogo

import (
	"testing"
)

func TestParserSearchForKeywords(t *testing.T) {
	psr := getTestParser("")
	psr.tokenizer.elements.insert(elementCategoryReleaseGroup, "THORA")
	for _, v := range *psr.tokenizer.tokens {
		v.Category = tokenCategoryUnknown
		if v.Content == "TaigaSubs" {
			v.Content = "THORA"
		}
	}
	psr.searchForKeywords()
	if psr.tokenizer.elements.AudioTerm[0] != "FLAC" {
		t.Errorf("expected \"FLAC\", got \"%s\"", psr.tokenizer.elements.AudioTerm[0])
	}
	if psr.tokenizer.elements.VideoResolution != "1280x720" {
		t.Errorf("expected \"1280x720\", got \"%s\"", psr.tokenizer.elements.VideoResolution)
	}
	if psr.tokenizer.elements.FileChecksum != "1234ABCD" {
		t.Errorf("expected \"1234ABCD\", got \"%s\"", psr.tokenizer.elements.FileChecksum)
	}
}

func TestParserSearchForIsolatedNumbers(t *testing.T) {
	psr := getTestParser("")
	psr.searchForIsolatedNumbers()
	if psr.tokenizer.elements.AnimeYear != "2008" {
		t.Errorf("expected \"2008\", got \"%s\"", psr.tokenizer.elements.AnimeYear)
	}
}

func TestParserSearchForEpisodeNumber(t *testing.T) {
	psr := getTestParser("")
	psr.searchForEpisodeNumber()
	if psr.tokenizer.elements.EpisodeNumber[0] != "01" {
		t.Errorf("expected \"01\", got \"%s\"", psr.tokenizer.elements.EpisodeNumber[0])
	}
	if psr.tokenizer.elements.ReleaseVersion[0] != "2" {
		t.Errorf("expected \"2\", got \"%s\"", psr.tokenizer.elements.ReleaseVersion[0])
	}
}

func TestParserSearchForAnimeTitle(t *testing.T) {
	psr := getTestParser("")
	psr.searchForKeywords()
	psr.searchForIsolatedNumbers()
	psr.searchForEpisodeNumber()
	psr.searchForAnimeTitle()
	if psr.tokenizer.elements.AnimeTitle != "Toradora!" {
		t.Errorf("expected \"Toradora!\", got \"%s\"", psr.tokenizer.elements.AnimeTitle)
	}
}

func TestParserSearchForReleaseGroup(t *testing.T) {
	psr := getTestParser("")
	psr.searchForReleaseGroup()
	if psr.tokenizer.elements.ReleaseGroup != "TaigaSubs" {
		t.Errorf("expected \"TaigaSubs\", got \"%s\"", psr.tokenizer.elements.ReleaseGroup)
	}
}

func TestParserSearchForEpisodeTitle(t *testing.T) {
	psr := getTestParser("")
	psr.searchForKeywords()
	psr.searchForIsolatedNumbers()
	psr.searchForEpisodeNumber()
	psr.searchForAnimeTitle()
	psr.searchForReleaseGroup()
	psr.searchForEpisodeTitle()
	if psr.tokenizer.elements.EpisodeTitle != "Tiger and Dragon" {
		t.Errorf("expected \"Tiger and Dragon\", got \"%s\"", psr.tokenizer.elements.EpisodeTitle)
	}
}
