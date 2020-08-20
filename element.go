package anitogo

type elementCategory int
type elements struct {
	AnimeSeason         []string `json:"anime_season,omitempty"`
	AnimeSeasonPrefix   []string `json:"anime_season_prefix,omitempty"`
	AnimeTitle          string   `json:"anime_title,omitempty"`
	AnimeType           string   `json:"anime_type,omitempty"`
	AnimeYear           string   `json:"anime_year,omitempty"`
	AudioTerm           []string `json:"audio_term,omitempty"`
	DeviceCompatibility []string `json:"device_compatibility,omitempty"`
	EpisodeNumber       []string `json:"episode_number,omitempty"`
	EpisodeNumberAlt    []string `json:"episode_number_alt,omitempty"`
	EpisodePrefix       []string `json:"episode_prefix,omitempty"`
	EpisodeTitle        string   `json:"episode_title,omitempty"`
	FileChecksum        string   `json:"file_checksum,omitempty"`
	FileExtension       string   `json:"file_extension,omitempty"`
	FileName            string   `json:"file_name,omitempty"`
	Language            []string `json:"language,omitempty"`
	Other               []string `json:"other,omitempty"`
	ReleaseGroup        string   `json:"release_group,omitempty"`
	ReleaseInformation  []string `json:"release_information,omitempty"`
	ReleaseVersion      []string `json:"release_version,omitempty"`
	Source              []string `json:"source,omitempty"`
	Subtitles           []string `json:"subtitles,omitempty"`
	VideoResolution     string   `json:"video_resolution,omitempty"`
	VideoTerm           []string `json:"video_term,omitempty"`
	VolumeNumber        []string `json:"volume_number,omitempty"`
	VolumePrefix        []string `json:"volume_prefix,omitempty"`
	Unknown             []string `json:"unknown,omitempty"`
	CheckAltNumber      bool     `json:"-"`
}

const (
	elementCategoryAnimeSeason elementCategory = iota
	elementCategoryAnimeSeasonPrefix
	elementCategoryAnimeTitle
	elementCategoryAnimeType
	elementCategoryAnimeYear
	elementCategoryAudioTerm
	elementCategoryDeviceCompatibility
	elementCategoryEpisodeNumber
	elementCategoryEpisodeNumberAlt
	elementCategoryEpisodePrefix
	elementCategoryEpisodeTitle
	elementCategoryFileChecksum
	elementCategoryFileExtension
	elementCategoryFileName
	elementCategoryLanguage
	elementCategoryOther
	elementCategoryReleaseGroup
	elementCategoryReleaseInformation
	elementCategoryReleaseVersion
	elementCategorySource
	elementCategorySubtitles
	elementCategoryVideoResolution
	elementCategoryVideoTerm
	elementCategoryVolumeNumber
	elementCategoryVolumePrefix
	elementCategoryUnknown
	elementCategoryCheckAltNumber
)

func newElements() *elements {
	e := elements{}

	return &e
}

func newElementsArr() []elements {
	e := []elements{}

	return e
}

func checkInList(arr []string, content string) bool {
	for _, v := range arr {
		if v == content {
			return true
		}
	}
	return false
}

func getIndex(arr []string, content string) int {
	for i, v := range arr {
		if v == content {
			return i
		}
	}
	return -1
}

func (e *elements) getCheckAltNumber() bool {
	return e.CheckAltNumber
}

func (e *elements) setCheckAltNumber(value bool) {
	e.CheckAltNumber = value
}

func (e *elements) insert(cat elementCategory, content string) {
	switch cat {
	case elementCategoryAnimeSeason:
		if checkInList(e.AnimeSeason, content) {
			return
		} else {
			e.AnimeSeason = append(e.AnimeSeason, content)
		}
	case elementCategoryAnimeSeasonPrefix:
		if checkInList(e.AnimeSeasonPrefix, content) {
			return
		} else {
			e.AnimeSeasonPrefix = append(e.AnimeSeasonPrefix, content)
		}
	case elementCategoryAnimeTitle:
		e.AnimeTitle = content
	case elementCategoryAnimeType:
		e.AnimeType = content
	case elementCategoryAnimeYear:
		e.AnimeYear = content
	case elementCategoryAudioTerm:
		if checkInList(e.AudioTerm, content) {
			return
		} else {
			e.AudioTerm = append(e.AudioTerm, content)
		}
	case elementCategoryDeviceCompatibility:
		if checkInList(e.DeviceCompatibility, content) {
			return
		} else {
			e.DeviceCompatibility = append(e.DeviceCompatibility, content)
		}
	case elementCategoryEpisodeNumber:
		if checkInList(e.EpisodeNumber, content) {
			return
		} else {
			e.EpisodeNumber = append(e.EpisodeNumber, content)
		}
	case elementCategoryEpisodeNumberAlt:
		if checkInList(e.EpisodeNumberAlt, content) {
			return
		} else {
			e.EpisodeNumberAlt = append(e.EpisodeNumberAlt, content)
		}
	case elementCategoryEpisodePrefix:
		if checkInList(e.EpisodePrefix, content) {
			return
		} else {
			e.EpisodePrefix = append(e.EpisodePrefix, content)
		}
	case elementCategoryEpisodeTitle:
		e.EpisodeTitle = content
	case elementCategoryFileChecksum:
		e.FileChecksum = content
	case elementCategoryFileExtension:
		e.FileExtension = content
	case elementCategoryFileName:
		e.FileName = content
	case elementCategoryLanguage:
		if checkInList(e.Language, content) {
			return
		} else {
			e.Language = append(e.Language, content)
		}
	case elementCategoryOther:
		if checkInList(e.Other, content) {
			return
		} else {
			e.Other = append(e.Other, content)
		}
	case elementCategoryReleaseGroup:
		e.ReleaseGroup = content
	case elementCategoryReleaseInformation:
		if checkInList(e.ReleaseInformation, content) {
			return
		} else {
			e.ReleaseInformation = append(e.ReleaseInformation, content)
		}
	case elementCategoryReleaseVersion:
		if checkInList(e.ReleaseVersion, content) {
			return
		} else {
			e.ReleaseVersion = append(e.ReleaseVersion, content)
		}
	case elementCategorySource:
		if checkInList(e.Source, content) {
			return
		} else {
			e.Source = append(e.Source, content)
		}
	case elementCategorySubtitles:
		if checkInList(e.Subtitles, content) {
			return
		} else {
			e.Subtitles = append(e.Subtitles, content)
		}
	case elementCategoryVideoResolution:
		e.VideoResolution = content
	case elementCategoryVideoTerm:
		if checkInList(e.VideoTerm, content) {
			return
		} else {
			e.VideoTerm = append(e.VideoTerm, content)
		}
	case elementCategoryVolumeNumber:
		if checkInList(e.VolumeNumber, content) {
			return
		} else {
			e.VolumeNumber = append(e.VolumeNumber, content)
		}
	case elementCategoryVolumePrefix:
		if checkInList(e.VolumePrefix, content) {
			return
		} else {
			e.VolumePrefix = append(e.VolumePrefix, content)
		}
	case elementCategoryUnknown:
		if checkInList(e.Unknown, content) {
			return
		} else {
			e.Unknown = append(e.Unknown, content)
		}
	}
}

func (e *elements) erase(cat elementCategory) {
	switch cat {
	case elementCategoryAnimeSeason:
		e.AnimeSeason = nil
	case elementCategoryAnimeSeasonPrefix:
		e.AnimeSeasonPrefix = nil
	case elementCategoryAnimeTitle:
		e.AnimeTitle = ""
	case elementCategoryAnimeType:
		e.AnimeType = ""
	case elementCategoryAnimeYear:
		e.AnimeYear = ""
	case elementCategoryAudioTerm:
		e.AudioTerm = nil
	case elementCategoryDeviceCompatibility:
		e.DeviceCompatibility = nil
	case elementCategoryEpisodeNumber:
		e.EpisodeNumber = nil
	case elementCategoryEpisodeNumberAlt:
		e.EpisodeNumberAlt = nil
	case elementCategoryEpisodePrefix:
		e.EpisodePrefix = nil
	case elementCategoryEpisodeTitle:
		e.EpisodeTitle = ""
	case elementCategoryFileChecksum:
		e.FileChecksum = ""
	case elementCategoryFileExtension:
		e.FileExtension = ""
	case elementCategoryFileName:
		e.FileName = ""
	case elementCategoryLanguage:
		e.Language = nil
	case elementCategoryOther:
		e.Other = nil
	case elementCategoryReleaseGroup:
		e.ReleaseGroup = ""
	case elementCategoryReleaseInformation:
		e.ReleaseInformation = nil
	case elementCategoryReleaseVersion:
		e.ReleaseVersion = nil
	case elementCategorySource:
		e.Source = nil
	case elementCategorySubtitles:
		e.Subtitles = nil
	case elementCategoryVideoResolution:
		e.VideoResolution = ""
	case elementCategoryVideoTerm:
		e.VideoTerm = nil
	case elementCategoryVolumeNumber:
		e.VolumeNumber = nil
	case elementCategoryVolumePrefix:
		e.VolumePrefix = nil
	case elementCategoryUnknown:
		e.Unknown = nil
	}
}

func (e *elements) fixJson() {
	for _, cat := range []elementCategory{
		elementCategoryAnimeSeason,
		elementCategoryAnimeSeasonPrefix,
		elementCategoryAudioTerm,
		elementCategoryDeviceCompatibility,
		elementCategoryEpisodeNumber,
		elementCategoryEpisodeNumberAlt,
		elementCategoryEpisodePrefix,
		elementCategoryLanguage,
		elementCategoryOther,
		elementCategoryReleaseInformation,
		elementCategorySubtitles,
		elementCategoryVideoTerm,
		elementCategoryVolumeNumber,
		elementCategoryVolumePrefix,
		elementCategoryUnknown,
	} {
		if equal(e.get(cat), []string{""}) {
			e.erase(cat)
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

func (e *elements) remove(cat elementCategory, content string) {
	switch cat {
	case elementCategoryAnimeSeason:
		idx := getIndex(e.AnimeSeason, content)
		if idx == -1 {
			return
		} else {
			e.AnimeSeason = append(e.AnimeSeason[:idx], e.AnimeSeason[idx+1:]...)
		}
	case elementCategoryAnimeSeasonPrefix:
		idx := getIndex(e.AnimeSeasonPrefix, content)
		if idx == -1 {
			return
		} else {
			e.AnimeSeasonPrefix = append(e.AnimeSeasonPrefix[:idx], e.AnimeSeasonPrefix[idx+1:]...)
		}
	case elementCategoryAnimeTitle:
		e.AnimeTitle = ""
	case elementCategoryAnimeType:
		e.AnimeType = ""
	case elementCategoryAnimeYear:
		e.AnimeYear = ""
	case elementCategoryAudioTerm:
		idx := getIndex(e.AudioTerm, content)
		if idx == -1 {
			return
		} else {
			e.AudioTerm = append(e.AudioTerm[:idx], e.AudioTerm[idx+1:]...)
		}
	case elementCategoryDeviceCompatibility:
		idx := getIndex(e.DeviceCompatibility, content)
		if idx == -1 {
			return
		} else {
			e.DeviceCompatibility = append(e.DeviceCompatibility[:idx], e.DeviceCompatibility[idx+1:]...)
		}
	case elementCategoryEpisodeNumber:
		idx := getIndex(e.EpisodeNumber, content)
		if idx == -1 {
			return
		} else {
			e.EpisodeNumber = append(e.EpisodeNumber[:idx], e.EpisodeNumber[idx+1:]...)
		}
	case elementCategoryEpisodeNumberAlt:
		idx := getIndex(e.EpisodeNumberAlt, content)
		if idx == -1 {
			return
		} else {
			e.EpisodeNumberAlt = append(e.EpisodeNumberAlt[:idx], e.EpisodeNumberAlt[idx+1:]...)
		}
	case elementCategoryEpisodePrefix:
		idx := getIndex(e.EpisodePrefix, content)
		if idx == -1 {
			return
		} else {
			e.EpisodePrefix = append(e.EpisodePrefix[:idx], e.EpisodePrefix[idx+1:]...)
		}
	case elementCategoryEpisodeTitle:
		e.EpisodeTitle = ""
	case elementCategoryFileChecksum:
		e.FileChecksum = ""
	case elementCategoryFileExtension:
		e.FileExtension = ""
	case elementCategoryFileName:
		e.FileName = ""
	case elementCategoryLanguage:
		idx := getIndex(e.Language, content)
		if idx == -1 {
			return
		} else {
			e.Language = append(e.Language[:idx], e.Language[idx+1:]...)
		}
	case elementCategoryOther:
		idx := getIndex(e.Other, content)
		if idx == -1 {
			return
		} else {
			e.Other = append(e.Other[:idx], e.Other[idx+1:]...)
		}
	case elementCategoryReleaseGroup:
		e.ReleaseGroup = ""
	case elementCategoryReleaseInformation:
		idx := getIndex(e.ReleaseInformation, content)
		if idx == -1 {
			return
		} else {
			e.ReleaseInformation = append(e.ReleaseInformation[:idx], e.ReleaseInformation[idx+1:]...)
		}
	case elementCategoryReleaseVersion:
		idx := getIndex(e.ReleaseVersion, content)
		if idx == -1 {
			return
		} else {
			e.ReleaseVersion = append(e.ReleaseVersion[:idx], e.ReleaseVersion[idx+1:]...)
		}
	case elementCategorySource:
		idx := getIndex(e.Source, content)
		if idx == -1 {
			return
		} else {
			e.Source = append(e.Source[:idx], e.Source[idx+1:]...)
		}
	case elementCategorySubtitles:
		idx := getIndex(e.Subtitles, content)
		if idx == -1 {
			return
		} else {
			e.Subtitles = append(e.Subtitles[:idx], e.Subtitles[idx+1:]...)
		}
	case elementCategoryVideoResolution:
		e.VideoResolution = ""
	case elementCategoryVideoTerm:
		idx := getIndex(e.VideoTerm, content)
		if idx == -1 {
			return
		} else {
			e.VideoTerm = append(e.VideoTerm[:idx], e.VideoTerm[idx+1:]...)
		}
	case elementCategoryVolumeNumber:
		idx := getIndex(e.VolumeNumber, content)
		if idx == -1 {
			return
		} else {
			e.VolumeNumber = append(e.VolumeNumber[:idx], e.VolumeNumber[idx+1:]...)
		}
	case elementCategoryVolumePrefix:
		idx := getIndex(e.VolumePrefix, content)
		if idx == -1 {
			return
		} else {
			e.VolumePrefix = append(e.VolumePrefix[:idx], e.VolumePrefix[idx+1:]...)
		}
	case elementCategoryUnknown:
		idx := getIndex(e.Unknown, content)
		if idx == -1 {
			return
		} else {
			e.Unknown = append(e.Unknown[:idx], e.Unknown[idx+1:]...)
		}
	}
}

func (e *elements) contains(cat elementCategory) bool {
	switch cat {
	case elementCategoryAnimeSeason:
		return e.AnimeSeason != nil
	case elementCategoryAnimeSeasonPrefix:
		return e.AnimeSeasonPrefix != nil
	case elementCategoryAnimeTitle:
		return e.AnimeTitle != ""
	case elementCategoryAnimeType:
		return e.AnimeType != ""
	case elementCategoryAnimeYear:
		return e.AnimeYear != ""
	case elementCategoryAudioTerm:
		return e.AudioTerm != nil
	case elementCategoryDeviceCompatibility:
		return e.DeviceCompatibility != nil
	case elementCategoryEpisodeNumber:
		return e.EpisodeNumber != nil
	case elementCategoryEpisodeNumberAlt:
		return e.EpisodeNumberAlt != nil
	case elementCategoryEpisodePrefix:
		return e.EpisodePrefix != nil
	case elementCategoryEpisodeTitle:
		return e.EpisodeTitle != ""
	case elementCategoryFileChecksum:
		return e.FileChecksum != ""
	case elementCategoryFileExtension:
		return e.FileExtension != ""
	case elementCategoryFileName:
		return e.FileName != ""
	case elementCategoryLanguage:
		return e.Language != nil
	case elementCategoryOther:
		return e.Other != nil
	case elementCategoryReleaseGroup:
		return e.ReleaseGroup != ""
	case elementCategoryReleaseInformation:
		return e.ReleaseInformation != nil
	case elementCategoryReleaseVersion:
		return e.ReleaseVersion != nil
	case elementCategorySource:
		return e.Source != nil
	case elementCategorySubtitles:
		return e.Subtitles != nil
	case elementCategoryVideoResolution:
		return e.VideoResolution != ""
	case elementCategoryVideoTerm:
		return e.VideoTerm != nil
	case elementCategoryVolumeNumber:
		return e.VolumeNumber != nil
	case elementCategoryVolumePrefix:
		return e.VolumePrefix != nil
	case elementCategoryUnknown:
		return e.Unknown != nil
	}
	return false
}

func (e *elements) empty() bool {
	return e == &elements{}
}

func (e *elements) get(cat elementCategory) []string {
	switch cat {
	case elementCategoryAnimeSeason:
		return e.AnimeSeason
	case elementCategoryAnimeSeasonPrefix:
		return e.AnimeSeasonPrefix
	case elementCategoryAnimeTitle:
		return []string{e.AnimeTitle}
	case elementCategoryAnimeType:
		return []string{e.AnimeType}
	case elementCategoryAnimeYear:
		return []string{e.AnimeYear}
	case elementCategoryAudioTerm:
		return e.AudioTerm
	case elementCategoryDeviceCompatibility:
		return e.DeviceCompatibility
	case elementCategoryEpisodeNumber:
		return e.EpisodeNumber
	case elementCategoryEpisodeNumberAlt:
		return e.EpisodeNumberAlt
	case elementCategoryEpisodePrefix:
		return e.EpisodePrefix
	case elementCategoryEpisodeTitle:
		return []string{e.EpisodeTitle}
	case elementCategoryFileChecksum:
		return []string{e.FileChecksum}
	case elementCategoryFileExtension:
		return []string{e.FileExtension}
	case elementCategoryFileName:
		return []string{e.FileName}
	case elementCategoryLanguage:
		return e.Language
	case elementCategoryOther:
		return e.Other
	case elementCategoryReleaseGroup:
		return []string{e.ReleaseGroup}
	case elementCategoryReleaseInformation:
		return e.ReleaseInformation
	case elementCategoryReleaseVersion:
		return e.ReleaseVersion
	case elementCategorySource:
		return e.Source
	case elementCategorySubtitles:
		return e.Subtitles
	case elementCategoryVideoResolution:
		return []string{e.VideoResolution}
	case elementCategoryVideoTerm:
		return e.VideoTerm
	case elementCategoryVolumeNumber:
		return e.VolumeNumber
	case elementCategoryVolumePrefix:
		return e.VolumePrefix
	case elementCategoryUnknown:
		return e.Unknown
	}
	return []string{}
}

func (e elementCategory) isSearchable() bool {
	searchableCategories := []elementCategory{
		elementCategoryAnimeSeasonPrefix,
		elementCategoryAnimeType,
		elementCategoryAudioTerm,
		elementCategoryDeviceCompatibility,
		elementCategoryEpisodePrefix,
		elementCategoryFileChecksum,
		elementCategoryLanguage,
		elementCategoryOther,
		elementCategoryReleaseGroup,
		elementCategoryReleaseInformation,
		elementCategoryReleaseVersion,
		elementCategorySource,
		elementCategorySubtitles,
		elementCategoryVideoResolution,
		elementCategoryVideoTerm,
		elementCategoryVolumePrefix,
	}

	for _, v := range searchableCategories {
		if e == v {
			return true
		}
	}
	return false
}

func (e elementCategory) isSingular() bool {
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
		if e == v {
			return false
		}
	}
	return true
}
