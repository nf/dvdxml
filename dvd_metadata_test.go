package main

import (
	"dvd-metadata-parser/dvd"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseDVDMetadata tests parsing of a real XML file
func TestParseDVDMetadata(t *testing.T) {
	// Use the first XML file for testing
	testFile := "source/s1d1.xml"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found, skipping test", testFile)
	}

	dvdData, err := dvd.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Failed to parse DVD metadata: %v", err)
	}

	// Basic validation
	if dvdData == nil {
		t.Fatal("DVD metadata is nil")
	}

	if dvdData.Device == "" {
		t.Error("Device should not be empty")
	}

	if len(dvdData.Tracks) == 0 {
		t.Error("Should have at least one track")
	}

	if dvdData.LongestTrack <= 0 {
		t.Error("Longest track should be greater than 0")
	}

	// Validate first track has expected fields
	if len(dvdData.Tracks) > 0 {
		track := dvdData.Tracks[0]
		if track.Index <= 0 {
			t.Error("Track index should be greater than 0")
		}
		if track.Length <= 0 {
			t.Error("Track length should be greater than 0")
		}
		if track.Width <= 0 || track.Height <= 0 {
			t.Error("Track should have valid resolution")
		}
		if track.Format == "" {
			t.Error("Track format should not be empty")
		}
		if track.FPS <= 0 {
			t.Error("Track FPS should be greater than 0")
		}

		// Validate audio streams
		for _, audio := range track.AudioStreams {
			if audio.Index <= 0 {
				t.Error("Audio stream index should be greater than 0")
			}
			if audio.Language == "" {
				t.Error("Audio stream language should not be empty")
			}
			if audio.Format == "" {
				t.Error("Audio stream format should not be empty")
			}
			if audio.Frequency <= 0 {
				t.Error("Audio stream frequency should be greater than 0")
			}
			if audio.Channels <= 0 {
				t.Error("Audio stream channels should be greater than 0")
			}
		}

		// Validate subtitle streams
		for _, sub := range track.SubtitleStreams {
			if sub.Index <= 0 {
				t.Error("Subtitle stream index should be greater than 0")
			}
			if sub.Language == "" {
				t.Error("Subtitle stream language should not be empty")
			}
		}

		// Validate chapters
		for _, chapter := range track.Chapters {
			if chapter.Index <= 0 {
				t.Error("Chapter index should be greater than 0")
			}
			if chapter.Length < 0 {
				t.Error("Chapter length should not be negative")
			}
			if chapter.StartCell <= 0 {
				t.Error("Chapter start cell should be greater than 0")
			}
		}
	}
}

// TestParseAllXMLFiles tests parsing all XML files in the source directory
func TestParseAllXMLFiles(t *testing.T) {
	sourceDir := "source"

	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		t.Skipf("Source directory %s not found, skipping test", sourceDir)
	}

	pattern := filepath.Join(sourceDir, "*.xml")
	xmlFiles, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("Error finding XML files: %v", err)
	}

	if len(xmlFiles) == 0 {
		t.Fatalf("No XML files found in %s", sourceDir)
	}

	successCount := 0
	for _, xmlFile := range xmlFiles {
		dvdData, err := dvd.ParseFile(xmlFile)
		if err != nil {
			t.Errorf("Failed to parse %s: %v", xmlFile, err)
			continue
		}

		// Basic validation
		if dvdData == nil {
			t.Errorf("DVD metadata is nil for file %s", xmlFile)
			continue
		}

		if len(dvdData.Tracks) == 0 {
			t.Errorf("No tracks found in file %s", xmlFile)
			continue
		}

		successCount++
	}

	t.Logf("Successfully parsed %d out of %d XML files", successCount, len(xmlFiles))

	if successCount == 0 {
		t.Fatal("Failed to parse any XML files")
	}
}

// TestSpecificFieldValues tests specific values from a known XML file
func TestSpecificFieldValues(t *testing.T) {
	testFile := "source/s1d1.xml"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found, skipping test", testFile)
	}

	dvdData, err := dvd.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Failed to parse DVD metadata: %v", err)
	}

	// Test known values from s1d1.xml
	if dvdData.Device != "./s1d1/Law And Order Svu" {
		t.Errorf("Expected device './s1d1/Law And Order Svu', got '%s'", dvdData.Device)
	}

	if dvdData.Title != "unknown" {
		t.Errorf("Expected title 'unknown', got '%s'", dvdData.Title)
	}

	if dvdData.VMGID != "DVDVIDEO-VMG" {
		t.Errorf("Expected VMG ID 'DVDVIDEO-VMG', got '%s'", dvdData.VMGID)
	}

	if dvdData.LongestTrack != 5 {
		t.Errorf("Expected longest track 5, got %d", dvdData.LongestTrack)
	}

	if len(dvdData.Tracks) != 10 {
		t.Errorf("Expected 10 tracks, got %d", len(dvdData.Tracks))
	}

	// Test first track specific values
	if len(dvdData.Tracks) > 0 {
		track := dvdData.Tracks[0]
		if track.Index != 1 {
			t.Errorf("Expected track index 1, got %d", track.Index)
		}
		if track.Width != 720 {
			t.Errorf("Expected width 720, got %d", track.Width)
		}
		if track.Height != 576 {
			t.Errorf("Expected height 576, got %d", track.Height)
		}
		if track.Format != "PAL" {
			t.Errorf("Expected format 'PAL', got '%s'", track.Format)
		}
		if track.FPS != 25.00 {
			t.Errorf("Expected FPS 25.00, got %.2f", track.FPS)
		}
		if track.Aspect != "4/3" {
			t.Errorf("Expected aspect '4/3', got '%s'", track.Aspect)
		}
		if len(track.AudioStreams) != 2 {
			t.Errorf("Expected 2 audio streams, got %d", len(track.AudioStreams))
		}
		if len(track.SubtitleStreams) != 4 {
			t.Errorf("Expected 4 subtitle streams, got %d", len(track.SubtitleStreams))
		}
		if len(track.Chapters) != 5 {
			t.Errorf("Expected 5 chapters, got %d", len(track.Chapters))
		}
	}
}

// TestDVDMethods tests the helper methods on the DVD struct
func TestDVDMethods(t *testing.T) {
	testFile := "source/s1d1.xml"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found, skipping test", testFile)
	}

	dvdData, err := dvd.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Failed to parse DVD metadata: %v", err)
	}

	// Test GetLongestTrack
	longestTrack := dvdData.GetLongestTrack()
	if longestTrack == nil {
		t.Error("GetLongestTrack should return a track")
	} else {
		if longestTrack.Index != 5 {
			t.Errorf("Expected longest track index 5, got %d", longestTrack.Index)
		}
	}

	// Test GetTrackByIndex
	track1 := dvdData.GetTrackByIndex(1)
	if track1 == nil {
		t.Error("GetTrackByIndex(1) should return a track")
	} else {
		if track1.Index != 1 {
			t.Errorf("Expected track index 1, got %d", track1.Index)
		}
	}

	// Test with non-existent track
	nonExistent := dvdData.GetTrackByIndex(999)
	if nonExistent != nil {
		t.Error("GetTrackByIndex(999) should return nil")
	}

	// Test GetTotalDuration
	totalDuration := dvdData.GetTotalDuration()
	if totalDuration <= 0 {
		t.Error("GetTotalDuration should return a positive value")
	}

	// Test GetAudioLanguages
	audioLangs := dvdData.GetAudioLanguages()
	if len(audioLangs) == 0 {
		t.Error("GetAudioLanguages should return at least one language")
	}

	// Test GetSubtitleLanguages
	subLangs := dvdData.GetSubtitleLanguages()
	if len(subLangs) == 0 {
		t.Error("GetSubtitleLanguages should return at least one language")
	}
}

// TestEpisodeMode tests the episode content finder
func TestEpisodeMode(t *testing.T) {
	testFile := "source/s1d1.xml"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found, skipping test", testFile)
	}

	dvdData, err := dvd.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Failed to parse DVD metadata: %v", err)
	}

	// Test the package method directly
	matches := dvdData.FindContentAroundDuration(40.0, 5.0)

	tracksFound := 0
	chaptersFound := 0

	for _, match := range matches {
		if match.Type == "track" {
			tracksFound++
		} else if match.Type == "chapter" {
			chaptersFound++
		}
	}

	// Based on s1d1.xml, we expect to find several tracks around 40 minutes
	if tracksFound == 0 {
		t.Error("Expected to find at least one track around 40 minutes")
	}

	t.Logf("Found %d tracks and %d chapters around 40 minutes", tracksFound, chaptersFound)
}

// isAroundDuration checks if a duration is within a tolerance of the target duration
func isAroundDuration(duration, target, tolerance float64) bool {
	return duration >= (target-tolerance) && duration <= (target+tolerance)
}

// TestIsAroundDuration tests the duration matching function
func TestIsAroundDuration(t *testing.T) {
	target := 40.0 * 60.0   // 40 minutes in seconds
	tolerance := 5.0 * 60.0 // 5 minutes tolerance

	testCases := []struct {
		duration    float64
		expected    bool
		description string
	}{
		{40.0 * 60.0, true, "exactly 40 minutes"},
		{35.0 * 60.0, true, "35 minutes (lower bound)"},
		{45.0 * 60.0, true, "45 minutes (upper bound)"},
		{41.5 * 60.0, true, "41.5 minutes (within range)"},
		{34.9 * 60.0, false, "34.9 minutes (below range)"},
		{45.1 * 60.0, false, "45.1 minutes (above range)"},
		{20.0 * 60.0, false, "20 minutes (well below)"},
		{60.0 * 60.0, false, "60 minutes (well above)"},
	}

	for _, tc := range testCases {
		result := isAroundDuration(tc.duration, target, tolerance)
		if result != tc.expected {
			t.Errorf("isAroundDuration(%.1f, %.1f, %.1f) = %v, expected %v for %s",
				tc.duration, target, tolerance, result, tc.expected, tc.description)
		}
	}
}

// TestFindEpisodeContentFunction tests the findEpisodeContent function behavior
func TestFindEpisodeContentFunction(t *testing.T) {
	testFile := "source/s1d1.xml"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found, skipping test", testFile)
	}

	dvdData, err := dvd.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Failed to parse DVD metadata: %v", err)
	}

	// Test with different durations
	testCases := []struct {
		targetMinutes    float64
		toleranceMinutes float64
		expectedMatches  int // approximate expected matches
		description      string
	}{
		{40.0, 5.0, 4, "40 minute episodes"}, // Should find several tracks
		{20.0, 2.0, 0, "20 minute episodes"}, // Should find none
		{164.0, 10.0, 1, "long content"},     // Should find the longest track
	}

	for _, tc := range testCases {
		matches := dvdData.FindContentAroundDuration(tc.targetMinutes, tc.toleranceMinutes)
		if tc.expectedMatches > 0 && len(matches) == 0 {
			t.Errorf("Expected to find matches for %s, got 0", tc.description)
		} else if tc.expectedMatches == 0 && len(matches) > 0 {
			t.Errorf("Expected no matches for %s, got %d", tc.description, len(matches))
		}
		t.Logf("%s: found %d matches", tc.description, len(matches))
	}
}

// TestFFmpegCommandGeneration tests the FFmpeg command generation
func TestFFmpegCommandGeneration(t *testing.T) {
	testFile := "source/s1d1.xml"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found, skipping test", testFile)
	}

	dvdData, err := dvd.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Failed to parse DVD metadata: %v", err)
	}

	// Test generateFFmpegCommand function
	matches := dvdData.FindContentAroundDuration(164.0, 10.0) // Should find the long track
	if len(matches) == 0 {
		t.Skip("No matches found for test")
	}

	match := matches[0]
	dvdPath := extractDVDPath(dvdData.Device)
	cmd := generateFFmpegCommand(match, dvdPath, "test_episodes")

	// Validate the command contains expected elements
	if !strings.Contains(cmd, "ffmpeg") {
		t.Error("Command should contain 'ffmpeg'")
	}
	if !strings.Contains(cmd, "dvdvideo:") {
		t.Error("Command should contain 'dvdvideo:'")
	}
	if !strings.Contains(cmd, "-c copy") {
		t.Error("Command should contain '-c copy' for stream copying")
	}
	if !strings.Contains(cmd, ".mkv") {
		t.Error("Command should output to .mkv file")
	}

	t.Logf("Generated FFmpeg command: %s", cmd)
}

// TestExtractDVDPath tests DVD path extraction from device strings
func TestExtractDVDPath(t *testing.T) {
	testCases := []struct {
		device   string
		expected string
	}{
		{"./s1d1/Law And Order Svu", "s1d1/Law And Order Svu"},
		{"s2d1/Some Movie", "s2d1/Some Movie"},
		{"/path/to/dvd", "/path/to/dvd"},
		{"./", ""},
	}

	for _, tc := range testCases {
		result := extractDVDPath(tc.device)
		if result != tc.expected {
			t.Errorf("extractDVDPath(%q) = %q, expected %q", tc.device, result, tc.expected)
		}
	}
}

// TestFFmpegModeOutput tests that FFmpeg mode only outputs commands
func TestFFmpegModeOutput(t *testing.T) {
	testFile := "source/s1d1.xml"

	// Check if test file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skipf("Test file %s not found, skipping test", testFile)
	}

	// This test would require capturing stdout, which is complex in Go tests
	// For now, we just validate that the logic paths work correctly
	dvdData, err := dvd.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Failed to parse DVD metadata: %v", err)
	}

	matches := dvdData.FindContentAroundDuration(164.0, 10.0) // Should find the long track
	if len(matches) == 0 {
		t.Skip("No matches found for test")
	}

	// Verify that FFmpeg commands are generated correctly
	for _, match := range matches {
		if match.Type == "track" {
			dvdPath := extractDVDPath(dvdData.Device)
			cmd := generateFFmpegCommand(match, dvdPath, "test")
			if !strings.Contains(cmd, "ffmpeg") {
				t.Error("Command should contain ffmpeg")
			}
			if !strings.Contains(cmd, "dvdvideo:") {
				t.Error("Command should use dvdvideo demuxer")
			}
		}
	}
}
