package anitogo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type failedParse struct {
	Expected Elements `json:"expected"`
	Got      Elements `json:"got"`
}

func TestParse(t *testing.T) {
	testDataPath := os.Getenv("TEST_DATA_PATH")
	if testDataPath == "" {
		t.Fatal("Missing TEST_DATA_PATH environment variable for json test data file")
	}

	e := []Elements{}
	notMatched := []failedParse{}
	jsonFile, err := os.Open(testDataPath)
	if err != nil {
		t.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	json.Unmarshal(byteValue, &e)
	for _, v := range e {
		ret, err := Parse(v.FileName, DefaultOptions)
		if err != nil {
			t.Error(err)
		}
		if !equal(v.AnimeSeason, ret.AnimeSeason) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.AnimeSeasonPrefix, ret.AnimeSeasonPrefix) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.AnimeTitle != v.AnimeTitle {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.AnimeYear != v.AnimeYear {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.AudioTerm, ret.AudioTerm) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.DeviceCompatibility, ret.DeviceCompatibility) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.EpisodeNumber, ret.EpisodeNumber) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.EpisodeNumberAlt, ret.EpisodeNumberAlt) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.EpisodePrefix, ret.EpisodePrefix) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.EpisodeTitle != v.EpisodeTitle {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.FileChecksum != v.FileChecksum {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.FileExtension != v.FileExtension {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.FileName != v.FileName {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.Language, ret.Language) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.ReleaseGroup != v.ReleaseGroup {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.ReleaseInformation, ret.ReleaseInformation) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.ReleaseVersion, ret.ReleaseVersion) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.Source, ret.Source) {
			if v.Source != nil {
				notMatched = append(notMatched, failedParse{
					Expected: v,
					Got:      *ret,
				})
			}
		} else if !equal(v.Subtitles, ret.Subtitles) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if ret.VideoResolution != v.VideoResolution {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.VideoTerm, ret.VideoTerm) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.VolumeNumber, ret.VolumeNumber) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		} else if !equal(v.VolumePrefix, ret.VolumePrefix) {
			notMatched = append(notMatched, failedParse{
				Expected: v,
				Got:      *ret,
			})
		}
	}

	if len(notMatched) > 0 {
		notMatchedJSON, err := json.MarshalIndent(notMatched, "", "    ")
		if err != nil {
			t.Error(err)
		}
		t.Fatalf("Failed %d/%d cases\n%s", len(notMatched), len(e), notMatchedJSON)
	}
}

func TestTraceError(t *testing.T) {
	err := traceError(indexTooLargeErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), indexTooLargeErr) == -1 {
			t.Errorf("expected %s in error, got %s", indexTooLargeErr, err.Error())
		}
	}
	err = traceError(indexTooSmallErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), indexTooSmallErr) == -1 {
			t.Errorf("expected %s in error, got %s", indexTooSmallErr, err.Error())
		}
	}
	err = traceError(endIndexTooSmallErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), endIndexTooSmallErr) == -1 {
			t.Errorf("expected %s in error, got %s", endIndexTooSmallErr, err.Error())
		}
	}
	err = traceError(tokensEmptyErr)
	if err == nil {
		t.Error("expected error got nil")
	} else {
		if strings.Index(err.Error(), tokensEmptyErr) == -1 {
			t.Errorf("expected %s in error, got %s", tokensEmptyErr, err.Error())
		}
	}
}

func BenchmarkParse(b *testing.B) {
	testDataPath := os.Getenv("TEST_DATA_PATH")
	if testDataPath == "" {
		b.Fatal("Missing TEST_DATA_PATH environment variable for json test data file")
	}
	e := []Elements{}
	jsonFile, err := os.Open(testDataPath)
	if err != nil {
		b.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		b.Fatal(err)
	}
	json.Unmarshal(byteValue, &e)
	for n := 0; n < b.N; n++ {
		for _, v := range e {
			_, err := Parse(v.FileName, DefaultOptions)
			if err != nil {
				b.Error(err)
			}
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
