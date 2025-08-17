// Package dvd provides DVD metadata parsing functionality for lsdvd XML output.
package dvd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// DVD represents the complete DVD metadata structure
type DVD struct {
	XMLName      xml.Name `xml:"lsdvd"`
	Device       string   `xml:"device"`
	Title        string   `xml:"title"`
	VMGID        string   `xml:"vmg_id"`
	ProviderID   string   `xml:"provider_id"`
	Tracks       []Track  `xml:"track"`
	LongestTrack int      `xml:"longest_track"`
}

// Track represents a DVD track with video, audio, subtitle, and chapter information
type Track struct {
	Index           int              `xml:"ix"`
	Length          float64          `xml:"length"`
	VTSID           string           `xml:"vts_id"`
	VTS             int              `xml:"vts"`
	TTN             int              `xml:"ttn"`
	FPS             float64          `xml:"fps"`
	Format          string           `xml:"format"`
	Aspect          string           `xml:"aspect"`
	Width           int              `xml:"width"`
	Height          int              `xml:"height"`
	DF              string           `xml:"df"`
	Palette         Palette          `xml:"palette"`
	Angles          int              `xml:"angles"`
	AudioStreams    []AudioStream    `xml:"audio"`
	SubtitleStreams []SubtitleStream `xml:"subp"`
	Chapters        []Chapter        `xml:"chapter"`
	Cells           []Cell           `xml:"cell"`
}

// Palette represents the color palette information
type Palette struct {
	Colors []string `xml:"color"`
}

// AudioStream represents an audio track
type AudioStream struct {
	Index        int    `xml:"ix"`
	LanguageCode string `xml:"langcode"`
	Language     string `xml:"language"`
	Format       string `xml:"format"`
	Frequency    int    `xml:"frequency"`
	Quantization string `xml:"quantization"`
	Channels     int    `xml:"channels"`
	APMode       int    `xml:"ap_mode"`
	Content      string `xml:"content"`
	StreamID     string `xml:"streamid"`
}

// SubtitleStream represents a subtitle track
type SubtitleStream struct {
	Index        int    `xml:"ix"`
	LanguageCode string `xml:"langcode"`
	Language     string `xml:"language"`
	Content      string `xml:"content"`
	StreamID     string `xml:"streamid"`
}

// Chapter represents a chapter within a track
type Chapter struct {
	Index     int     `xml:"ix"`
	Length    float64 `xml:"length"`
	StartCell int     `xml:"startcell"`
}

// Cell represents a cell within a track
type Cell struct {
	Index  int     `xml:"ix"`
	Length float64 `xml:"length"`
}

// ParseFile parses a single XML file and returns DVD metadata
func ParseFile(filename string) (*DVD, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	return ParseBytes(data)
}

// ParseBytes parses DVD metadata from XML byte data
func ParseBytes(data []byte) (*DVD, error) {
	// Fix common XML entity issues in the data
	// Fix malformed entity &Scan -> &amp;Scan
	data = bytes.ReplaceAll(data, []byte("Pan&Scan"), []byte("Pan&amp;Scan"))
	// Fix other potential malformed entities
	data = bytes.ReplaceAll(data, []byte("&Letterbox"), []byte("&amp;Letterbox"))

	var dvd DVD
	err := xml.Unmarshal(data, &dvd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	return &dvd, nil
}

// GetLongestTrack returns the longest track from the DVD, or nil if not found
func (d *DVD) GetLongestTrack() *Track {
	if d.LongestTrack > 0 && d.LongestTrack <= len(d.Tracks) {
		return &d.Tracks[d.LongestTrack-1] // Convert to 0-based index
	}
	return nil
}

// GetTrackByIndex returns a track by its index (1-based), or nil if not found
func (d *DVD) GetTrackByIndex(index int) *Track {
	for i := range d.Tracks {
		if d.Tracks[i].Index == index {
			return &d.Tracks[i]
		}
	}
	return nil
}

// GetTotalDuration returns the total duration of all tracks in seconds
func (d *DVD) GetTotalDuration() float64 {
	var total float64
	for _, track := range d.Tracks {
		total += track.Length
	}
	return total
}

// GetAudioLanguages returns a slice of unique audio languages across all tracks
func (d *DVD) GetAudioLanguages() []string {
	languageMap := make(map[string]bool)
	for _, track := range d.Tracks {
		for _, audio := range track.AudioStreams {
			if audio.Language != "" {
				languageMap[audio.Language] = true
			}
		}
	}

	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}
	return languages
}

// GetSubtitleLanguages returns a slice of unique subtitle languages across all tracks
func (d *DVD) GetSubtitleLanguages() []string {
	languageMap := make(map[string]bool)
	for _, track := range d.Tracks {
		for _, sub := range track.SubtitleStreams {
			if sub.Language != "" {
				languageMap[sub.Language] = true
			}
		}
	}

	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}
	return languages
}

// ContentMatch represents a track or chapter that matches certain criteria
type ContentMatch struct {
	Type     string   // "track" or "chapter"
	Track    *Track   // The track containing this content
	Chapter  *Chapter // The chapter (nil if Type is "track")
	Duration float64  // Duration in seconds
}

// FindContentAroundDuration finds tracks and chapters with duration around the target
func (d *DVD) FindContentAroundDuration(targetMinutes, toleranceMinutes float64) []ContentMatch {
	targetSeconds := targetMinutes * 60.0
	toleranceSeconds := toleranceMinutes * 60.0

	var matches []ContentMatch

	for i := range d.Tracks {
		track := &d.Tracks[i]

		// Check if the entire track matches
		if track.Length >= (targetSeconds-toleranceSeconds) && track.Length <= (targetSeconds+toleranceSeconds) {
			matches = append(matches, ContentMatch{
				Type:     "track",
				Track:    track,
				Chapter:  nil,
				Duration: track.Length,
			})
			continue // Don't check chapters if the whole track matches
		}

		// Check chapters within this track
		for j := range track.Chapters {
			chapter := &track.Chapters[j]
			if chapter.Length >= (targetSeconds-toleranceSeconds) && chapter.Length <= (targetSeconds+toleranceSeconds) {
				matches = append(matches, ContentMatch{
					Type:     "chapter",
					Track:    track,
					Chapter:  chapter,
					Duration: chapter.Length,
				})
			}
		}
	}

	return matches
}

// FindFortyMinuteContent is a convenience method to find content around 40 minutes
func (d *DVD) FindFortyMinuteContent() []ContentMatch {
	return d.FindContentAroundDuration(40.0, 5.0)
}
