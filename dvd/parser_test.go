package dvd

import (
	"testing"
)

// TestParseBytes tests parsing from byte data
func TestParseBytes(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<lsdvd>
    <device>./test</device>
    <title>Test DVD</title>
    <vmg_id>DVDVIDEO-VMG</vmg_id>
    <provider_id>TEST</provider_id>
    <track>
        <ix>1</ix>
        <length>100.0</length>
        <vts_id>DVDVIDEO-VTS</vts_id>
        <vts>1</vts>
        <ttn>1</ttn>
        <fps>25.00</fps>
        <format>PAL</format>
        <aspect>4/3</aspect>
        <width>720</width>
        <height>576</height>
        <df>?</df>
        <angles>1</angles>
        <audio>
            <ix>1</ix>
            <langcode>en</langcode>
            <language>English</language>
            <format>ac3</format>
            <frequency>48000</frequency>
            <channels>2</channels>
        </audio>
        <chapter>
            <ix>1</ix>
            <length>100.0</length>
            <startcell>1</startcell>
        </chapter>
    </track>
    <longest_track>1</longest_track>
</lsdvd>`)

	dvd, err := ParseBytes(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	if dvd.Device != "./test" {
		t.Errorf("Expected device './test', got '%s'", dvd.Device)
	}

	if dvd.Title != "Test DVD" {
		t.Errorf("Expected title 'Test DVD', got '%s'", dvd.Title)
	}

	if len(dvd.Tracks) != 1 {
		t.Errorf("Expected 1 track, got %d", len(dvd.Tracks))
	}

	if dvd.LongestTrack != 1 {
		t.Errorf("Expected longest track 1, got %d", dvd.LongestTrack)
	}
}

// TestEntityFixes tests the XML entity fixing functionality
func TestEntityFixes(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<lsdvd>
    <device>./test</device>
    <title>Test DVD</title>
    <vmg_id>DVDVIDEO-VMG</vmg_id>
    <provider_id>TEST</provider_id>
    <track>
        <ix>1</ix>
        <length>100.0</length>
        <df>Pan&Scan</df>
        <format>PAL</format>
        <width>720</width>
        <height>576</height>
    </track>
    <longest_track>1</longest_track>
</lsdvd>`)

	dvd, err := ParseBytes(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse XML with entity issues: %v", err)
	}

	if len(dvd.Tracks) != 1 {
		t.Errorf("Expected 1 track, got %d", len(dvd.Tracks))
	}

	if dvd.Tracks[0].DF != "Pan&Scan" {
		t.Errorf("Expected DF 'Pan&Scan', got '%s'", dvd.Tracks[0].DF)
	}
}

// TestDVDMethods tests helper methods on DVD struct
func TestDVDMethods(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<lsdvd>
    <device>./test</device>
    <title>Test DVD</title>
    <vmg_id>DVDVIDEO-VMG</vmg_id>
    <provider_id>TEST</provider_id>
    <track>
        <ix>1</ix>
        <length>100.0</length>
        <format>PAL</format>
        <audio>
            <ix>1</ix>
            <language>English</language>
        </audio>
        <subp>
            <ix>1</ix>
            <language>Spanish</language>
        </subp>
    </track>
    <track>
        <ix>2</ix>
        <length>200.0</length>
        <format>NTSC</format>
        <audio>
            <ix>1</ix>
            <language>French</language>
        </audio>
    </track>
    <longest_track>2</longest_track>
</lsdvd>`)

	dvd, err := ParseBytes(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	// Test GetLongestTrack
	longestTrack := dvd.GetLongestTrack()
	if longestTrack == nil {
		t.Fatal("GetLongestTrack should return a track")
	}
	if longestTrack.Index != 2 {
		t.Errorf("Expected longest track index 2, got %d", longestTrack.Index)
	}

	// Test GetTrackByIndex
	track1 := dvd.GetTrackByIndex(1)
	if track1 == nil {
		t.Fatal("GetTrackByIndex(1) should return a track")
	}
	if track1.Index != 1 {
		t.Errorf("Expected track index 1, got %d", track1.Index)
	}

	nonExistent := dvd.GetTrackByIndex(999)
	if nonExistent != nil {
		t.Error("GetTrackByIndex(999) should return nil")
	}

	// Test GetTotalDuration
	totalDuration := dvd.GetTotalDuration()
	expectedTotal := 300.0 // 100.0 + 200.0
	if totalDuration != expectedTotal {
		t.Errorf("Expected total duration %.1f, got %.1f", expectedTotal, totalDuration)
	}

	// Test GetAudioLanguages
	audioLangs := dvd.GetAudioLanguages()
	if len(audioLangs) != 2 {
		t.Errorf("Expected 2 audio languages, got %d", len(audioLangs))
	}

	// Test GetSubtitleLanguages
	subLangs := dvd.GetSubtitleLanguages()
	if len(subLangs) != 1 {
		t.Errorf("Expected 1 subtitle language, got %d", len(subLangs))
	}
}

// TestInvalidXML tests error handling for invalid XML
func TestInvalidXML(t *testing.T) {
	invalidXML := []byte(`<invalid>xml</incomplete>`)

	_, err := ParseBytes(invalidXML)
	if err == nil {
		t.Error("Expected error for invalid XML, got nil")
	}
}

// TestFindContentAroundDuration tests the content finding functionality
func TestFindContentAroundDuration(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<lsdvd>
    <device>./test</device>
    <title>Test DVD</title>
    <vmg_id>DVDVIDEO-VMG</vmg_id>
    <provider_id>TEST</provider_id>
    <track>
        <ix>1</ix>
        <length>2400.0</length>
        <format>PAL</format>
        <width>720</width>
        <height>576</height>
        <chapter>
            <ix>1</ix>
            <length>1200.0</length>
            <startcell>1</startcell>
        </chapter>
        <chapter>
            <ix>2</ix>
            <length>2400.0</length>
            <startcell>2</startcell>
        </chapter>
    </track>
    <track>
        <ix>2</ix>
        <length>3600.0</length>
        <format>PAL</format>
        <width>720</width>
        <height>576</height>
        <chapter>
            <ix>1</ix>
            <length>2400.0</length>
            <startcell>1</startcell>
        </chapter>
    </track>
    <longest_track>2</longest_track>
</lsdvd>`)

	dvd, err := ParseBytes(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	// Test finding content around 40 minutes (2400 seconds)
	matches := dvd.FindContentAroundDuration(40.0, 5.0)

	// Should find track 1 (2400s) and chapter 1 in track 2 (2400s)
	// Note: track 1 matches so its chapters are not checked
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}

	foundTrack := false
	foundChapter := false

	for _, match := range matches {
		if match.Type == "track" && match.Track.Index == 1 {
			foundTrack = true
			if match.Duration != 2400.0 {
				t.Errorf("Expected track duration 2400.0, got %.1f", match.Duration)
			}
		} else if match.Type == "chapter" && match.Track.Index == 2 && match.Chapter.Index == 1 {
			foundChapter = true
			if match.Duration != 2400.0 {
				t.Errorf("Expected chapter duration 2400.0, got %.1f", match.Duration)
			}
		}
	}

	if !foundTrack {
		t.Error("Expected to find track 1 in matches")
	}
	if !foundChapter {
		t.Error("Expected to find chapter 1 of track 2 in matches")
	}
}

// TestFindFortyMinuteContent tests the convenience method
func TestFindFortyMinuteContent(t *testing.T) {
	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<lsdvd>
    <device>./test</device>
    <title>Test DVD</title>
    <vmg_id>DVDVIDEO-VMG</vmg_id>
    <provider_id>TEST</provider_id>
    <track>
        <ix>1</ix>
        <length>2400.0</length>
        <format>PAL</format>
    </track>
    <track>
        <ix>2</ix>
        <length>1800.0</length>
        <format>PAL</format>
    </track>
    <longest_track>1</longest_track>
</lsdvd>`)

	dvd, err := ParseBytes(xmlData)
	if err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}

	// Test the convenience method
	matches := dvd.FindFortyMinuteContent()

	// Should find track 1 (2400s = 40 minutes)
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}

	if len(matches) > 0 {
		match := matches[0]
		if match.Type != "track" {
			t.Errorf("Expected match type 'track', got '%s'", match.Type)
		}
		if match.Track.Index != 1 {
			t.Errorf("Expected track index 1, got %d", match.Track.Index)
		}
		if match.Duration != 2400.0 {
			t.Errorf("Expected duration 2400.0, got %.1f", match.Duration)
		}
	}
}
