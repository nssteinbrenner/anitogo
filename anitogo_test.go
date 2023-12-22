package anitogo

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

type failedParse struct {
	Expected Elements `json:"expected"`
	Got      Elements `json:"got"`
}

var testDataPath = flag.String("file", "./test/data.json", "Path to test data JSON")

func TestAnitogoParse(t *testing.T) {
	retElems := Parse("", DefaultOptions)
	if retElems.FileName != "" {
		t.Error("expected empty elements")
	}
	retElems = Parse("1", DefaultOptions)
	if retElems.AnimeTitle != "" {
		t.Error("expected empty anime title")
	}
	noReleaseGroup := DefaultOptions
	noReleaseGroup.ParseReleaseGroup = false
	retElems = Parse("[THORA]_Toradora!_(2008)_-_01v2_-_Tiger_and_Dragon_[1280x720_H.264_FLAC][1234ABCD].mkv", noReleaseGroup)
	if retElems.ReleaseGroup != "" {
		t.Error("expected empty release group")
	}

	e := []Elements{}
	notMatched := []failedParse{}
	jsonFile, err := os.Open(*testDataPath)
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
		ret := Parse(v.FileName, DefaultOptions)
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
		} else if !equal(v.AnimeType, ret.AnimeType) {
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

func TestAnitogoRemoveIgnoredStrings(t *testing.T) {
	s := removeIgnoredStrings("testing this", []string{" ", "this"})
	if s != "testing" {
		t.Errorf("expected \"testing\" got \"%s\"", s)
	}
}

func TestAnitogoRemoveExtensionFromFilename(t *testing.T) {
	s := "[HorribleSubs] Boku no Hero Academia - 01 [1080p].mkv"
	kwm := newKeywordManager()
	filename, extension := removeExtensionFromFilename(kwm, s)
	if filename != "[HorribleSubs] Boku no Hero Academia - 01 [1080p]" {
		t.Errorf("expected \"[HorribleSubs] Boku no Hero Academia - 01 [1080p]\", got \"%s\"", filename)
	}
	if extension != "mkv" {
		t.Errorf("expected \"mkv\", got \"%s\"", extension)
	}
}

func BenchmarkAnitogoParse(b *testing.B) {
	e := []Elements{}
	jsonFile, err := os.Open(*testDataPath)
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
			Parse(v.FileName, DefaultOptions)
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
