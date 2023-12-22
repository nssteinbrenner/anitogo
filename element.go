package anitogo

type elementCategory int

// Elements is a struct representing a parsed anime filename.
type Elements struct {
	// Slice of strings representing the season of anime. "S1-S3" would be represented as []string{"1", "3"}.
	AnimeSeason []string `json:"anime_season,omitempty"`

	// Represents the strings prefixing the season in the file, e.g in "SEASON 2" "SEASON" is the AnimeSeasonPrefix.
	AnimeSeasonPrefix []string `json:"anime_season_prefix,omitempty"`

	// Title of the Anime. e.g in "[HorribleSubs] Boku no Hero Academia - 01 [1080p].mkv",
	// "Boku No Hero Academia" is the AnimeTitle.
	AnimeTitle string `json:"anime_title,omitempty"`

	// Type specified in the anime file, e.g ED, OP, Movie, etc.
	AnimeType string `json:"anime_type,omitempty"`

	// Year the anime was released.
	AnimeYear string `json:"anime_year,omitempty"`

	// Slice of strings representing the audio terms included in the filename, e.g FLAC, AAC, etc.
	AudioTerm []string `json:"audio_term,omitempty"`

	// Slice of strings representing devices the video is compatible with that are mentioned in the filename.
	DeviceCompatibility []string `json:"device_compatibility,omitempty"`

	// Slice of strings representing the episode numbers. "01-10" would be respresented as []string{"1", "10"}.
	EpisodeNumber []string `json:"episode_number,omitempty"`

	// Slice of strings representing the alternative episode number.
	// This is for cases where you may have an episode number relative to the season,
	// but a larger episode number as if it were all one season.
	// e.g in [Hatsuyuki]_Kuroko_no_Basuke_S3_-_01_(51)_[720p][10bit][619C57A0].mkv
	// 01 would be the EpisodeNumber, and 51 would be the EpisodeNumberAlt.
	EpisodeNumberAlt []string `json:"episode_number_alt,omitempty"`

	// Slice of strings representing the words prefixing the episode number in the file, e.g in "EPISODE 2", "EPISODE" is the prefix.
	EpisodePrefix []string `json:"episode_prefix,omitempty"`

	// Title of the episode. e.g in "[BM&T] Toradora! - 07v2 - Pool Opening [720p Hi10 ] [BD] [8F59F2BA]",
	// "Pool Opening" is the EpisodeTitle.
	EpisodeTitle string `json:"episode_title,omitempty"`

	// Checksum of the file, in [BM&T] Toradora! - 07v2 - Pool Opening [720p Hi10 ] [BD] [8F59F2BA],
	// "8F59F2BA" would be the FileChecksum.
	FileChecksum string `json:"file_checksum,omitempty"`

	// File extension, in [HorribleSubs] Boku no Hero Academia - 01 [1080p].mkv,
	// "mkv" would be the FileExtension.
	FileExtension string `json:"file_extension,omitempty"`

	// Full filename that was parsed.
	FileName string `json:"file_name,omitempty"`

	// Languages specified in the file name, e.g RU, JP, EN etc.
	Language []string `json:"language,omitempty"`

	// Terms that could not be parsed into other buckets, but were deemed identifiers.
	// In [chibi-Doki] Seikon no Qwaser - 13v0 (Uncensored Director's Cut) [988DB090].mkv,
	// "Uncensored" is parsed into Other.
	Other []string `json:"other,omitempty"`

	// The fan sub group that uploaded the file. In [HorribleSubs] Boku no Hero Academia - 01 [1080p],
	// "HorribleSubs" is the ReleaseGroup.
	ReleaseGroup string `json:"release_group,omitempty"`

	// Information about the release that wasn't a version.
	// In "[SubDESU-H] Swing out Sisters Complete Version (720p x264 8bit AC3) [3ABD57E6].mp4
	// "Complete" is parsed into ReleaseInformation.
	ReleaseInformation []string `json:"release_information,omitempty"`

	// Slice of strings representing the version of the release.
	// In [FBI] Baby Princess 3D Paradise Love 01v0 [BD][720p-AAC][457CC066].mkv, 0 is parsed into ReleaseVersion.
	ReleaseVersion []string `json:"release_version,omitempty"`

	// Slice of strings representing where the video was ripped from. e.g BLU-RAY, DVD, etc.
	Source []string `json:"source,omitempty"`

	// Slice of strings representing the type of subtitles included, e.g HARDSUB, BIG5, etc.
	Subtitles []string `json:"subtitles,omitempty"`

	// Resolution of the video. Can be formatted like 1920x1080, 1080, 1080p, etc depending
	// on how it is represented in the filename.
	VideoResolution string `json:"video_resolution,omitempty"`

	// Slice of strings representing the video terms included in the filename, e.g h264, x264, etc.
	VideoTerm []string `json:"video_term,omitempty"`

	// Slice of strings represnting the volume numbers. "01-10" would be represented as []string{"1", "10"}.
	VolumeNumber []string `json:"volume_number,omitempty"`

	// Slice of strings representing the words prefixing the volume number in the file, e.g in "VOLUME 2", "VOLUME" is the prefix.
	VolumePrefix []string `json:"volume_prefix,omitempty"`

	// Entries that could not be parsed into any other categories.
	Unknown []string `json:"unknown,omitempty"`

	// Bool determining if "EpisodeNumberAlt" should be parsed or not.
	CheckAltNumber bool `json:"-"`
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
)

func newElements() *Elements {
	e := Elements{}

	return &e
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

func (e *Elements) getCheckAltNumber() bool {
	return e.CheckAltNumber
}

func (e *Elements) setCheckAltNumber(value bool) {
	e.CheckAltNumber = value
}

func (e *Elements) insert(cat elementCategory, content string) {
	switch cat {
	case elementCategoryAnimeSeason:
		if checkInList(e.AnimeSeason, content) {
			return
		}
		e.AnimeSeason = append(e.AnimeSeason, content)
	case elementCategoryAnimeSeasonPrefix:
		if checkInList(e.AnimeSeasonPrefix, content) {
			return
		}
		e.AnimeSeasonPrefix = append(e.AnimeSeasonPrefix, content)
	case elementCategoryAnimeTitle:
		e.AnimeTitle = content
	case elementCategoryAnimeType:
		e.AnimeType = content
	case elementCategoryAnimeYear:
		e.AnimeYear = content
	case elementCategoryAudioTerm:
		if checkInList(e.AudioTerm, content) {
			return
		}
		e.AudioTerm = append(e.AudioTerm, content)
	case elementCategoryDeviceCompatibility:
		if checkInList(e.DeviceCompatibility, content) {
			return
		}
		e.DeviceCompatibility = append(e.DeviceCompatibility, content)
	case elementCategoryEpisodeNumber:
		if checkInList(e.EpisodeNumber, content) {
			return
		}
		e.EpisodeNumber = append(e.EpisodeNumber, content)
	case elementCategoryEpisodeNumberAlt:
		if checkInList(e.EpisodeNumberAlt, content) {
			return
		}
		e.EpisodeNumberAlt = append(e.EpisodeNumberAlt, content)
	case elementCategoryEpisodePrefix:
		if checkInList(e.EpisodePrefix, content) {
			return
		}
		e.EpisodePrefix = append(e.EpisodePrefix, content)
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
		}
		e.Language = append(e.Language, content)
	case elementCategoryOther:
		if checkInList(e.Other, content) {
			return
		}
		e.Other = append(e.Other, content)
	case elementCategoryReleaseGroup:
		e.ReleaseGroup = content
	case elementCategoryReleaseInformation:
		if checkInList(e.ReleaseInformation, content) {
			return
		}
		e.ReleaseInformation = append(e.ReleaseInformation, content)
	case elementCategoryReleaseVersion:
		if checkInList(e.ReleaseVersion, content) {
			return
		}
		e.ReleaseVersion = append(e.ReleaseVersion, content)
	case elementCategorySource:
		if checkInList(e.Source, content) {
			return
		}
		e.Source = append(e.Source, content)
	case elementCategorySubtitles:
		if checkInList(e.Subtitles, content) {
			return
		}
		e.Subtitles = append(e.Subtitles, content)
	case elementCategoryVideoResolution:
		e.VideoResolution = content
	case elementCategoryVideoTerm:
		if checkInList(e.VideoTerm, content) {
			return
		}
		e.VideoTerm = append(e.VideoTerm, content)
	case elementCategoryVolumeNumber:
		if checkInList(e.VolumeNumber, content) {
			return
		}
		e.VolumeNumber = append(e.VolumeNumber, content)
	case elementCategoryVolumePrefix:
		if checkInList(e.VolumePrefix, content) {
			return
		}
		e.VolumePrefix = append(e.VolumePrefix, content)
	case elementCategoryUnknown:
		if checkInList(e.Unknown, content) {
			return
		}
		e.Unknown = append(e.Unknown, content)
	}
}

func (e *Elements) erase(cat elementCategory) {
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

func (e *Elements) remove(cat elementCategory, content string) {
	switch cat {
	case elementCategoryAnimeSeason:
		idx := getIndex(e.AnimeSeason, content)
		if idx == -1 {
			return
		}
		e.AnimeSeason = append(e.AnimeSeason[:idx], e.AnimeSeason[idx+1:]...)
	case elementCategoryAnimeSeasonPrefix:
		idx := getIndex(e.AnimeSeasonPrefix, content)
		if idx == -1 {
			return
		}
		e.AnimeSeasonPrefix = append(e.AnimeSeasonPrefix[:idx], e.AnimeSeasonPrefix[idx+1:]...)
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
		}
		e.AudioTerm = append(e.AudioTerm[:idx], e.AudioTerm[idx+1:]...)
	case elementCategoryDeviceCompatibility:
		idx := getIndex(e.DeviceCompatibility, content)
		if idx == -1 {
			return
		}
		e.DeviceCompatibility = append(e.DeviceCompatibility[:idx], e.DeviceCompatibility[idx+1:]...)
	case elementCategoryEpisodeNumber:
		idx := getIndex(e.EpisodeNumber, content)
		if idx == -1 {
			return
		}
		e.EpisodeNumber = append(e.EpisodeNumber[:idx], e.EpisodeNumber[idx+1:]...)
	case elementCategoryEpisodeNumberAlt:
		idx := getIndex(e.EpisodeNumberAlt, content)
		if idx == -1 {
			return
		}
		e.EpisodeNumberAlt = append(e.EpisodeNumberAlt[:idx], e.EpisodeNumberAlt[idx+1:]...)
	case elementCategoryEpisodePrefix:
		idx := getIndex(e.EpisodePrefix, content)
		if idx == -1 {
			return
		}
		e.EpisodePrefix = append(e.EpisodePrefix[:idx], e.EpisodePrefix[idx+1:]...)
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
		}
		e.Language = append(e.Language[:idx], e.Language[idx+1:]...)
	case elementCategoryOther:
		idx := getIndex(e.Other, content)
		if idx == -1 {
			return
		}
		e.Other = append(e.Other[:idx], e.Other[idx+1:]...)
	case elementCategoryReleaseGroup:
		e.ReleaseGroup = ""
	case elementCategoryReleaseInformation:
		idx := getIndex(e.ReleaseInformation, content)
		if idx == -1 {
			return
		}
		e.ReleaseInformation = append(e.ReleaseInformation[:idx], e.ReleaseInformation[idx+1:]...)
	case elementCategoryReleaseVersion:
		idx := getIndex(e.ReleaseVersion, content)
		if idx == -1 {
			return
		}
		e.ReleaseVersion = append(e.ReleaseVersion[:idx], e.ReleaseVersion[idx+1:]...)
	case elementCategorySource:
		idx := getIndex(e.Source, content)
		if idx == -1 {
			return
		}
		e.Source = append(e.Source[:idx], e.Source[idx+1:]...)
	case elementCategorySubtitles:
		idx := getIndex(e.Subtitles, content)
		if idx == -1 {
			return
		}
		e.Subtitles = append(e.Subtitles[:idx], e.Subtitles[idx+1:]...)
	case elementCategoryVideoResolution:
		e.VideoResolution = ""
	case elementCategoryVideoTerm:
		idx := getIndex(e.VideoTerm, content)
		if idx == -1 {
			return
		}
		e.VideoTerm = append(e.VideoTerm[:idx], e.VideoTerm[idx+1:]...)
	case elementCategoryVolumeNumber:
		idx := getIndex(e.VolumeNumber, content)
		if idx == -1 {
			return
		}
		e.VolumeNumber = append(e.VolumeNumber[:idx], e.VolumeNumber[idx+1:]...)
	case elementCategoryVolumePrefix:
		idx := getIndex(e.VolumePrefix, content)
		if idx == -1 {
			return
		}
		e.VolumePrefix = append(e.VolumePrefix[:idx], e.VolumePrefix[idx+1:]...)
	case elementCategoryUnknown:
		idx := getIndex(e.Unknown, content)
		if idx == -1 {
			return
		}
		e.Unknown = append(e.Unknown[:idx], e.Unknown[idx+1:]...)
	}
}

func (e *Elements) contains(cat elementCategory) bool {
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

func (e *Elements) get(cat elementCategory) []string {
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
