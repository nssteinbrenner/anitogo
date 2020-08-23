# Anitogo [![Go Report Card](https://goreportcard.com/badge/github.com/nssteinbrenner/anitogo)](https://goreportcard.com/report/github.com/nssteinbrenner/anitogo) [![Build Status](https://travis-ci.org/nssteinbrenner/anitogo.svg?branch=master)](https://travis-ci.org/nssteinbrenner/anitogo) [![Coverage Status](https://coveralls.io/repos/github/nssteinbrenner/anitogo/badge.svg?branch=master)](https://coveralls.io/github/nssteinbrenner/anitogo?branch=master)
Anitogo is a Golang library for parsing anime video filenames. It is based off of [Anitomy](https://github.com/erengy/anitomy) and [Anitopy](https://github.com/igorcmoura/anitopy).

## Example
The following filename...

    [TaigaSubs]_Toradora!_(2008)_-_01v2_-_Tiger_and_Dragon_[1280x720_H.264_FLAC][1234ABCD].mkv

...is resolved into these elements:

- Release Group: "TaigaSubs"
- Anime Title: "Toradora!"
- Anime Year: "2008"
- Episode Number: ["01"]
- Release version: ["2"]
- Episode title: "Tiger and Dragon"
- Video resolution: "1280x720"
- Video term: ["H.264"]
- Audio term: ["FLAC"]
- File checksum: "1234ABCD"

The following example code:

    package main

    import (
        "fmt"
        "encoding/json"

        "github.com/nssteinbrenner/anitogo"
    )

    func main() {
        parsed, err := anitogo.Parse("[Nubles] Space Battleship Yamato 2199 (2012) episode 18 (720p 10 bit AAC)[1F56D642]", anitogo.DefaultOptions)
        if err != nil {
            fmt.Println(err)
        }
        jsonParsed, err := json.MarshalIndent(parsed, "", "    ")
        if err != nil {
            fmt.Println(err)
        }
        fmt.Println(string(jsonParsed) + "\n")

        // Accessing the elements directly
        fmt.Println("Anime Title:", parsed.AnimeTitle)
        fmt.Println("Anime Year:", parsed.AnimeYear)
        fmt.Println("Episode Number:", parsed.EpisodeNumber)
        fmt.Println("Release Group:", parsed.ReleaseGroup)
        fmt.Println("File Checksum:", parsed.FileChecksum)
    }

Will output:

    {
        "anime_title": "Space Battleship Yamato 2199",
        "anime_year": "2012",
        "audio_term": [
            "AAC"
        ],
        "episode_number": [
            "18"
        ],
        "file_checksum": "1F56D642",
        "file_name": "[Nubles] Space Battleship Yamato 2199 (2012) episode 18 (720p 10 bit AAC)[1F56D642]",
        "release_group": "Nubles",
        "video_resolution": "720p"
    }

    Anime Title: Space Battleship Yamato 2199
    Anime Year: 2012
    Episode Number: [18]
    Release Group: Nubles
    File Checksum: 1F56D642

The Parse function returns an element struct and an error. The full definition of the struct is here:

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
## Installation
Get the package:

    go get github.com/nssteinbrenner/anitogo

Then, import it in your code:

    import "github.com/nssteinbrenner/anitogo"

## Options
The Parse function receives the filename and an Options struct. The default options are as follows:

    var DefaultOptions = Options{
        AllowedDelimiters:  " _.&+,|", // Parse these as delimiters
        IgnoredStrings:     []string{}, // Ignore these when they are in the filename
        ParseEpisodeNumber: true, // Parse the episode number and include it in the elements
        ParseEpisodeTitle:  true, // Parse the episode title and include it in the elements
        ParseFileExtension: true, // Parse the file extension and include it in the elements
        ParseReleaseGroup:  true, // Parse the release group and include it in the elements
    }
