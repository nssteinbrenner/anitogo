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
	kws := make(map[string]keyword)
	kwfileExtensions := make(map[string]keyword)
	kwm := &keywordManager{
		keywords:       kws,
		fileExtensions: kwfileExtensions,
	}

	kwm.Add(elementCategoryAnimeSeasonPrefix, keywordOptionsUnidentifiable, []string{"S", "SAISON", "SEASON"})
	kwm.Add(elementCategoryAnimeType, keywordOptionsUnidentifiable, []string{
		"GEKIJOUBAN", "MOVIE", "OAD", "OAV", "ONA", "OVA", "SPECIAL", "SPECIALS", "TV"})
	kwm.Add(elementCategoryAnimeType, keywordOptionsUnidentifiableUnsearchable, []string{
		"SP"})
	kwm.Add(elementCategoryAnimeType, keywordOptionsUnidentifiableInvalid, []string{
		"ED", "ENDING", "NCED", "NCOP", "OP", "OPENING", "PREVIEW", "PV"})
	kwm.Add(elementCategoryAudioTerm, keywordOptionsDefault, []string{
		"2.0CH", "2CH", "5.1", "5.1CH", "DTS", "DTS-ES", "DTS5.1", "TRUEHD5.1",
		"AAC", "AACX2", "AACX3", "AACX4", "AC3", "EAC3", "E-AC-3", "FLAC",
		"FLACX2", "FLACX3", "FLACX4", "LOSSLESS", "MP3", "OGG", "VORBIS",
		"DD2", "DD2.0", "DUALAUDIO", "DUAL AUDIO"})
	kwm.Add(elementCategoryDeviceCompatibility, keywordOptionsDefault, []string{
		"IPAD3", "IPHONE5", "IPOD", "PS3", "XBOX", "XBOX360"})
	kwm.Add(elementCategoryDeviceCompatibility, keywordOptionsUnidentifiable, []string{
		"ANDROID"})
	kwm.Add(elementCategoryEpisodePrefix, keywordOptionsDefault, []string{
		"EP", "EP.", "EPS", "EPS.", "EPISODE", "EPISODE.", "EPISODES",
		"CAPITULO", "EPISODIO", "FOLGE"})
	kwm.Add(elementCategoryEpisodePrefix, keywordOptionsInvalid, []string{
		"E", "\x7B2C"})
	kwm.Add(elementCategoryFileExtension, keywordOptionsDefault, []string{
		"3GP", "AVI", "DIVX", "FLV", "M2TS", "MKV", "MOV", "MP4", "MPG",
		"OGM", "RM", "RMVB", "TS", "WEBM", "WMV"})
	kwm.Add(elementCategoryFileExtension, keywordOptionsInvalid, []string{
		"AAC", "AIFF", "FLAC", "M4A", "MP3", "MKA", "OGG", "WAV", "WMA",
		"7Z", "RAR", "ZIP", "ASS", "SRT"})
	kwm.Add(elementCategoryLanguage, keywordOptionsDefault, []string{
		"ENG", "ENGLISH", "ESPANOL", "JAP", "PT-BR", "SPANISH", "VOSTFR"})
	kwm.Add(elementCategoryLanguage, keywordOptionsUnidentifiable, []string{
		"ESP", "ITA"})
	kwm.Add(elementCategoryOther, keywordOptionsDefault, []string{
		"REMASTER", "REMASTERED", "UNCENSORED", "UNCUT", "TS", "VFR",
		"WIDESCREEN", "WS"})
	kwm.Add(elementCategoryReleaseGroup, keywordOptionsDefault, []string{
		"THORA", "HORRIBLESUBS", "ERAI-RAWS"})
	kwm.Add(elementCategoryReleaseInformation, keywordOptionsDefault, []string{
		"BATCH", "COMPLETE", "PATCH", "REMUX"})
	kwm.Add(elementCategoryReleaseInformation, keywordOptionsUnidentifiable, []string{
		"END", "FINAL"})
	kwm.Add(elementCategoryReleaseVersion, keywordOptionsDefault, []string{
		"V0", "V1", "V2", "V3", "V4"})
	kwm.Add(elementCategorySource, keywordOptionsDefault, []string{
		"BD", "BDRIP", "BLURAY", "BLU-RAY", "DVD", "DVD5", "DVD9",
		"DVD-R2J", "DVDRIP", "DVD-RIP", "R2DVD", "R2J", "R2JDVD",
		"R2JDVDRIP", "HDTV", "HDTVRIP", "TVRIP", "TV-RIP",
		"WEBCAST", "WEBRIP"})
	kwm.Add(elementCategorySubtitles, keywordOptionsDefault, []string{
		"ASS", "BIG5", "DUB", "DUBBED", "HARDSUB", "HARDSUBS", "RAW",
		"SOFTSUB", "SOFTSUBS", "SUB", "SUBBED", "SUBTITLED"})
	kwm.Add(elementCategoryVideoTerm, keywordOptionsDefault, []string{
		"23.976FPS", "24FPS", "29.97FPS", "30FPS", "60FPS", "120FPS",
		"8BIT", "8-BIT", "10BIT", "10BITS", "10-BIT", "10-BITS",
		"HI10", "HI10P", "HI444", "HI444P", "HI444PP",
		"H264", "H265", "H.264", "H.265", "X264", "X265", "X.264",
		"AVC", "HEVC", "HEVC2", "DIVX", "DIVX5", "DIVX6", "XVID",
		"AVI", "RMVB", "WMV", "WMV3", "WMV9",
		"HQ", "LQ",
		"HD", "SD"})
	kwm.Add(elementCategoryVolumePrefix, keywordOptionsDefault, []string{
		"VOL", "VOL.", "VOLUME"})

	return kwm
}

func (kd keyword) empty() bool {
	return kd == keyword{}
}

func (kwm *keywordManager) Add(cat elementCategory, opt keywordOption, keywords []string) {
	for _, kw := range keywords {
		if len(kw) == 0 {
			continue
		}
		if cat != elementCategoryFileExtension {
			v := kwm.keywords[kw]
			if !v.empty() {
				continue
			}
			kwm.keywords[kw] = keyword{
				Category: cat,
				Options:  opt,
			}
		} else {
			v := kwm.fileExtensions[kw]
			if !v.empty() {
				continue
			}
			kwm.fileExtensions[kw] = keyword{
				Category: cat,
				Options:  opt,
			}
		}
	}
}

func (kwm *keywordManager) find(word string, cat elementCategory) (keyword, bool) {
	var kw keyword

	if word == "" {
		return kw, false
	}

	if cat != elementCategoryFileExtension {
		v := kwm.keywords[word]
		if !v.empty() {
			kw = v
			if kw.Category != elementCategoryUnknown && kw.Category != cat {
				return kw, false
			}
			return kw, true
		}
	} else {
		v := kwm.fileExtensions[word]
		if !v.empty() {
			kw = v
			if kw.Category != elementCategoryUnknown && kw.Category != cat {
				return kw, false
			}
			return kw, true
		}
	}
	return kw, false
}

func (kwm *keywordManager) FindWithoutCategory(word string) (keyword, bool) {
	var kw keyword

	if word == "" {
		return kw, false
	}
	v := kwm.keywords[word]
	if !v.empty() {
		kw = v
		return kw, true
	}
	v = kwm.keywords[word]
	if !v.empty() {
		kw = v
		return kw, true
	}
	return kw, false
}

func (kwm *keywordManager) Peek(word string, e *Elements) (indexSets, error) {
	entries := map[elementCategory][]string{
		elementCategoryAudioTerm:       []string{"Dual Audio", "DualAudio"},
		elementCategoryVideoTerm:       []string{"H264", "H.264", "h264", "h.264"},
		elementCategoryVideoResolution: []string{"480p", "720p", "1080p"},
		elementCategorySource:          []string{"Blu-Ray"},
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

func (kwm *keywordManager) Normalize(text string) string {
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
