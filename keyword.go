package anitogo

import (
	"sort"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type indexSet struct {
	BeginPos int
	EndPos   int
}

type indexSets []indexSet

type keywordOption struct {
	Identifiable bool
	Searchable   bool
	Valid        bool
}

type keyword struct {
	Category elementCategory
	Options  keywordOption
}

type keywordManager struct {
	keywords       map[string]keyword
	fileExtensions map[string]keyword
}

var (
	keywordOptionsDefault = keywordOption{
		Identifiable: true,
		Searchable:   true,
		Valid:        true,
	}
	keywordOptionsInvalid = keywordOption{
		Identifiable: true,
		Searchable:   true,
		Valid:        false,
	}
	keywordOptionsUnidentifiable = keywordOption{
		Identifiable: false,
		Searchable:   true,
		Valid:        true,
	}
	keywordOptionsUnidentifiableInvalid = keywordOption{
		Identifiable: false,
		Searchable:   true,
		Valid:        false,
	}
	keywordOptionsUnidentifiableUnsearchable = keywordOption{
		Identifiable: false,
		Searchable:   false,
		Valid:        true,
	}
)

func newKeywordManager() *keywordManager {
	kwm := &keywordManager{
		keywords:       make(map[string]keyword),
		fileExtensions: make(map[string]keyword),
	}

	kwm.add(elementCategoryAnimeSeasonPrefix, keywordOptionsUnidentifiable, []string{"S", "SAISON", "SEASON"})
	kwm.add(elementCategoryAnimeType, keywordOptionsUnidentifiable, []string{
		"GEKIJOUBAN", "MOVIE", "OAD", "OAV", "ONA", "OVA", "SPECIAL", "SPECIALS", "TV"})
	kwm.add(elementCategoryAnimeType, keywordOptionsUnidentifiableUnsearchable, []string{
		"SP"}) // e.g "Yumeiro Patissiere SP Professional"
	kwm.add(elementCategoryAnimeType, keywordOptionsUnidentifiableInvalid, []string{
		"ED", "ENDING", "NCED", "NCOP", "OP", "OPENING", "PREVIEW", "PV"})
	kwm.add(elementCategoryAudioTerm, keywordOptionsDefault, []string{
		// Audio channels
		"2.0CH", "2CH", "5.1", "5.1CH", "DTS", "DTS-ES", "DTS5.1", "TRUEHD5.1",
		// Audio codec
		"AAC", "AACX2", "AACX3", "AACX4", "AC3", "EAC3", "E-AC-3", "FLAC",
		"FLACX2", "FLACX3", "FLACX4", "LOSSLESS", "MP3", "OGG", "VORBIS",
		"DD2", "DD2.0",
		// Audio language
		"DUALAUDIO", "DUAL AUDIO"})
	kwm.add(elementCategoryDeviceCompatibility, keywordOptionsDefault, []string{
		"IPAD3", "IPHONE5", "IPOD", "PS3", "XBOX", "XBOX360"})
	kwm.add(elementCategoryDeviceCompatibility, keywordOptionsUnidentifiable, []string{
		"ANDROID"})
	kwm.add(elementCategoryEpisodePrefix, keywordOptionsDefault, []string{
		"EP", "EP.", "EPS", "EPS.", "EPISODE", "EPISODE.", "EPISODES",
		"CAPITULO", "EPISODIO", "FOLGE"})
	kwm.add(elementCategoryEpisodePrefix, keywordOptionsInvalid, []string{
		"E", "\x7B2C"}) // Single letter episode keywords are not valid tokens
	kwm.add(elementCategoryFileExtension, keywordOptionsDefault, []string{
		"3GP", "AVI", "DIVX", "FLV", "M2TS", "MKV", "MOV", "MP4", "MPG",
		"OGM", "RM", "RMVB", "TS", "WEBM", "WMV"})
	kwm.add(elementCategoryFileExtension, keywordOptionsInvalid, []string{
		"AAC", "AIFF", "FLAC", "M4A", "MP3", "MKA", "OGG", "WAV", "WMA",
		"7Z", "RAR", "ZIP", "ASS", "SRT"})
	kwm.add(elementCategoryLanguage, keywordOptionsDefault, []string{
		"ENG", "ENGLISH", "ESPANOL", "JAP", "PT-BR", "SPANISH", "VOSTFR"})
	kwm.add(elementCategoryLanguage, keywordOptionsUnidentifiable, []string{
		"ESP", "ITA"}) // e.g "Tokyo ESP", "Bokura ga Ita"
	kwm.add(elementCategoryOther, keywordOptionsDefault, []string{
		"REMASTER", "REMASTERED", "UNCENSORED", "UNCUT", "TS", "VFR",
		"WIDESCREEN", "WS"})
	kwm.add(elementCategoryReleaseGroup, keywordOptionsDefault, []string{
		"THORA", "HORRIBLESUBS", "ERAI-RAWS"})
	kwm.add(elementCategoryReleaseInformation, keywordOptionsDefault, []string{
		"BATCH", "COMPLETE", "PATCH", "REMUX"})
	kwm.add(elementCategoryReleaseInformation, keywordOptionsUnidentifiable, []string{
		"END", "FINAL"}) // e.g "The End of Evangelion", "Final Approach"
	kwm.add(elementCategoryReleaseVersion, keywordOptionsDefault, []string{
		"V0", "V1", "V2", "V3", "V4"})
	kwm.add(elementCategorySource, keywordOptionsDefault, []string{
		"BD", "BDRIP", "BLURAY", "BLU-RAY", "DVD", "DVD5", "DVD9",
		"DVD-R2J", "DVDRIP", "DVD-RIP", "R2DVD", "R2J", "R2JDVD",
		"R2JDVDRIP", "HDTV", "HDTVRIP", "TVRIP", "TV-RIP",
		"WEBCAST", "WEBRIP"})
	kwm.add(elementCategorySubtitles, keywordOptionsDefault, []string{
		"ASS", "BIG5", "DUB", "DUBBED", "HARDSUB", "HARDSUBS", "RAW",
		"SOFTSUB", "SOFTSUBS", "SUB", "SUBBED", "SUBTITLED"})
	kwm.add(elementCategoryVideoTerm, keywordOptionsDefault, []string{
		// Frame rate
		"23.976FPS", "24FPS", "29.97FPS", "30FPS", "60FPS", "120FPS",
		// Video codec
		"8BIT", "8-BIT", "10BIT", "10BITS", "10-BIT", "10-BITS",
		"HI10", "HI10P", "HI444", "HI444P", "HI444PP",
		"H264", "H265", "H.264", "H.265", "X264", "X265", "X.264",
		"AVC", "HEVC", "HEVC2", "DIVX", "DIVX5", "DIVX6", "XVID",
		// Video format
		"AVI", "RMVB", "WMV", "WMV3", "WMV9",
		// Video quality
		"HQ", "LQ",
		// Video resolution
		"HD", "SD"})
	kwm.add(elementCategoryVolumePrefix, keywordOptionsDefault, []string{
		"VOL", "VOL.", "VOLUME"})

	return kwm
}

func (kd keyword) empty() bool {
	return kd == keyword{}
}

func (kwm *keywordManager) add(cat elementCategory, opt keywordOption, keywords []string) {
	for _, kw := range keywords {
		if len(kw) == 0 {
			continue
		}
		if cat != elementCategoryFileExtension {
			kwm.keywords[kw] = keyword{
				Category: cat,
				Options:  opt,
			}
		} else {
			kwm.fileExtensions[kw] = keyword{
				Category: cat,
				Options:  opt,
			}
		}
	}
}

func (kwm *keywordManager) find(word string, cat elementCategory) (keyword, bool) {
	if word == "" {
		return keyword{}, false
	}

	if cat != elementCategoryFileExtension {
		v := kwm.keywords[word]
		if !v.empty() && (v.Category == elementCategoryUnknown || v.Category == cat) {
			return v, true
		}
	} else {
		v := kwm.fileExtensions[word]
		if !v.empty() {
			return v, true
		}
	}
	return keyword{}, false
}

func (kwm *keywordManager) findWithoutCategory(word string) (keyword, bool) {
	if word == "" {
		return keyword{}, false
	}
	v := kwm.keywords[word]
	if !v.empty() {
		return v, true
	}
	v = kwm.fileExtensions[word]
	if !v.empty() {
		return v, true
	}
	return keyword{}, false
}

func (kwm *keywordManager) peek(word string, e *Elements) (indexSets, error) {
	entries := map[elementCategory][]string{
		elementCategoryAudioTerm:       {"Dual Audio", "DualAudio"},
		elementCategoryVideoTerm:       {"H264", "H.264", "h264", "h.264"},
		elementCategoryVideoResolution: {"480p", "720p", "1080p"},
		elementCategorySource:          {"Blu-Ray"},
	}

	preIdentifiedTokens := indexSets{}

	for cat, keywords := range entries {
		for _, kw := range keywords {
			keywordBeginPos := strings.Index(word, kw)
			if keywordBeginPos != -1 {
				e.insert(cat, kw)
				keywordEndPos := keywordBeginPos + len(kw)
				if keywordEndPos > len(word) {
					return indexSets{}, traceError(indexTooLargeErr)
				}
				preIdentifiedTokens = append(preIdentifiedTokens, indexSet{keywordBeginPos, keywordEndPos})
			}
		}
	}
	sort.Sort(preIdentifiedTokens)
	return preIdentifiedTokens, nil
}

func (kwm *keywordManager) normalize(text string) string {
	f := norm.Form(3)

	return strings.ToUpper(string(f.Bytes([]byte(text))))
}

func (idxSet indexSets) Len() int {
	return len(idxSet)
}

func (idxSet indexSets) Less(i, j int) bool {
	return (idxSet[i].BeginPos + idxSet[i].EndPos) < (idxSet[j].BeginPos + idxSet[j].EndPos)
}

func (idxSet indexSets) Swap(i, j int) {
	idxSet[i], idxSet[j] = idxSet[j], idxSet[i]
}
