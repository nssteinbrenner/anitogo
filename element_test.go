package anitogo

import (
	"testing"
)

var multiElementFields = []elementCategory{
	elementCategoryAnimeSeason,
	elementCategoryAnimeSeasonPrefix,
	elementCategoryAnimeType,
	elementCategoryAudioTerm,
	elementCategoryDeviceCompatibility,
	elementCategoryEpisodeNumber,
	elementCategoryEpisodeNumberAlt,
	elementCategoryEpisodePrefix,
	elementCategoryLanguage,
	elementCategoryOther,
	elementCategoryReleaseInformation,
	elementCategoryReleaseVersion,
	elementCategorySource,
	elementCategorySubtitles,
	elementCategoryVideoTerm,
	elementCategoryVolumeNumber,
	elementCategoryVolumePrefix,
	elementCategoryUnknown,
}

var singleElementFields = []elementCategory{
	elementCategoryAnimeTitle,
	elementCategoryAnimeYear,
	elementCategoryEpisodeTitle,
	elementCategoryFileChecksum,
	elementCategoryFileExtension,
	elementCategoryFileName,
	elementCategoryReleaseGroup,
	elementCategoryVideoResolution,
}

func TestElementMultiElementFields(t *testing.T) {
	e := &Elements{}
	for _, v := range multiElementFields {
		found, _ := e.getMultiElementField(v)
		if found != true {
			t.Error("expected true, got false")
		}
		e.insert(v, "test")
		_, field := e.getMultiElementField(v)
		if len(*field) != 1 {
			t.Errorf("expected len == 1 got %d", len(*field))
		}
		e.erase(v)
		if len(*field) != 0 {
			t.Errorf("expected len == 0 got %d", len(*field))
		}
		e.insert(v, "test")
		e.remove(v, "test")
		if len(*field) != 0 {
			t.Errorf("expected len == 0 got %d", len(*field))
		}
		e.erase(v)
		found = e.contains(v)
		if found != false {
			t.Error("expected false, got true")
		}
	}
	found, _ := e.getMultiElementField(100)
	if found != false {
		t.Error("expected false, got true")
	}

	for _, v := range singleElementFields {
		found, _ := e.getMultiElementField(v)
		if found != false {
			t.Error("expected false, got true")
		}
	}

}

func TestElementSingleElementFields(t *testing.T) {
	e := &Elements{}
	for _, v := range singleElementFields {
		found, _ := e.getSingleElementField(v)
		if found != true {
			t.Error("expected true, got false")
		}
		e.insert(v, "test")
		_, field := e.getSingleElementField(v)
		if *field != "test" {
			t.Errorf("expected \"test\" got \"%s\"", *field)
		}
		e.erase(v)
		if *field != "" {
			t.Errorf("expected \"\" got \"%s\"", *field)
		}
		e.insert(v, "test")
		e.remove(v, "test")
		if *field != "" {
			t.Errorf("expected \"\" got \"%s\"", *field)
		}
		e.erase(v)
		found = e.contains(v)
		if found != false {
			t.Error("expected false, got true")
		}
	}
	found, _ := e.getSingleElementField(100)
	if found != false {
		t.Error("expected false, got true")
	}

	for _, v := range multiElementFields {
		found, _ := e.getSingleElementField(v)
		if found != false {
			t.Error("expected false, got true")
		}
	}
}

func TestElementIsSearchable(t *testing.T) {
	nonSearchableCategories := []elementCategory{
		elementCategoryAnimeSeason,
		elementCategoryAnimeTitle,
		elementCategoryEpisodeNumber,
		elementCategoryEpisodeNumberAlt,
		elementCategoryEpisodeTitle,
		elementCategoryFileName,
		elementCategoryVolumeNumber,
		elementCategoryUnknown,
	}

	for _, v := range nonSearchableCategories {
		i := v.isSearchable()
		if i {
			t.Errorf("expected false, got true")
		}
	}
}

func TestElementIsSingular(t *testing.T) {
	nonSingularCategories := []elementCategory{
		elementCategoryAnimeSeason,
		elementCategoryAnimeType,
		elementCategoryAudioTerm,
		elementCategoryDeviceCompatibility,
		elementCategoryEpisodeNumber,
		elementCategoryLanguage,
		elementCategoryOther,
		elementCategoryReleaseInformation,
		elementCategorySource,
		elementCategoryVideoTerm,
	}

	for _, v := range nonSingularCategories {
		i := v.isSingular()
		if i {
			t.Errorf("expected false, got true")
		}
	}
}

func TestElementCheckInList(t *testing.T) {
	arr := []string{"test", "test1"}
	i := checkInList(arr, "test2")
	if i {
		t.Errorf("expected false, got true")
	}
	i = checkInList(arr, "test1")
	if !i {
		t.Errorf("expected true, got false")
	}
}

func TestElementGetIndex(t *testing.T) {
	arr := []string{"test", "test1"}
	i := getIndex(arr, "test2")
	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}

	i = getIndex(arr, "test1")
	if i != 1 {
		t.Errorf("expected 1, got %d", i)
	}
}

func TestElementRemove(t *testing.T) {
	e := &Elements{}
	e.remove(elementCategoryEpisodeNumber, "1A")
}
