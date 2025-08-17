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
