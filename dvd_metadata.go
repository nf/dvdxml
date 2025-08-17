package main

import (
	"dvd-metadata-parser/dvd"
	"fmt"
	"os"
	"path/filepath"
)

// printDVDSummary prints a summary of the DVD metadata
func printDVDSummary(filename string, dvdData *dvd.DVD) {
	fmt.Printf("\n=== %s ===\n", filename)
	fmt.Printf("Device: %s\n", dvdData.Device)
	fmt.Printf("Title: %s\n", dvdData.Title)
	fmt.Printf("Provider ID: %s\n", dvdData.ProviderID)
	fmt.Printf("Number of tracks: %d\n", len(dvdData.Tracks))
	fmt.Printf("Longest track: %d\n", dvdData.LongestTrack)

	for i, track := range dvdData.Tracks {
		fmt.Printf("\n  Track %d:\n", track.Index)
		fmt.Printf("    Length: %.2f seconds (%.2f minutes)\n", track.Length, track.Length/60)
		fmt.Printf("    Resolution: %dx%d\n", track.Width, track.Height)
		fmt.Printf("    Aspect: %s\n", track.Aspect)
		fmt.Printf("    Format: %s @ %.2f fps\n", track.Format, track.FPS)
		fmt.Printf("    Chapters: %d\n", len(track.Chapters))
		fmt.Printf("    Audio streams: %d\n", len(track.AudioStreams))
		fmt.Printf("    Subtitle streams: %d\n", len(track.SubtitleStreams))

		// Show audio stream details
		for j, audio := range track.AudioStreams {
			fmt.Printf("      Audio %d: %s (%s) - %s, %d Hz, %d channels\n",
				audio.Index, audio.Language, audio.LanguageCode,
				audio.Format, audio.Frequency, audio.Channels)
			if j >= 2 { // Limit output for readability
				if len(track.AudioStreams) > 3 {
					fmt.Printf("      ... and %d more audio streams\n", len(track.AudioStreams)-3)
				}
				break
			}
		}

		// Show subtitle stream details
		for j, sub := range track.SubtitleStreams {
			fmt.Printf("      Subtitle %d: %s (%s)\n",
				sub.Index, sub.Language, sub.LanguageCode)
			if j >= 2 { // Limit output for readability
				if len(track.SubtitleStreams) > 3 {
					fmt.Printf("      ... and %d more subtitle streams\n", len(track.SubtitleStreams)-3)
				}
				break
			}
		}

		if i >= 4 { // Limit number of tracks shown for readability
			if len(dvdData.Tracks) > 5 {
				fmt.Printf("\n  ... and %d more tracks\n", len(dvdData.Tracks)-5)
			}
			break
		}
	}
}

// printDetailedTrackInfo prints detailed information about a specific track
func printDetailedTrackInfo(track dvd.Track) {
	fmt.Printf("\n--- Detailed Track %d Information ---\n", track.Index)
	fmt.Printf("Length: %.2f seconds\n", track.Length)
	fmt.Printf("Video: %s, %dx%d, %s, %.2f fps\n", track.Format, track.Width, track.Height, track.Aspect, track.FPS)
	fmt.Printf("VTS: %d, TTN: %d\n", track.VTS, track.TTN)

	fmt.Printf("\nAudio Streams (%d):\n", len(track.AudioStreams))
	for _, audio := range track.AudioStreams {
		fmt.Printf("  [%d] %s (%s): %s, %d Hz, %d ch, Stream ID: %s\n",
			audio.Index, audio.Language, audio.LanguageCode,
			audio.Format, audio.Frequency, audio.Channels, audio.StreamID)
	}

	fmt.Printf("\nSubtitle Streams (%d):\n", len(track.SubtitleStreams))
	for _, sub := range track.SubtitleStreams {
		fmt.Printf("  [%d] %s (%s): %s, Stream ID: %s\n",
			sub.Index, sub.Language, sub.LanguageCode, sub.Content, sub.StreamID)
	}

	fmt.Printf("\nChapters (%d):\n", len(track.Chapters))
	for _, chapter := range track.Chapters {
		fmt.Printf("  Chapter %d: %.2f seconds (starts at cell %d)\n",
			chapter.Index, chapter.Length, chapter.StartCell)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dvd_metadata.go <source_directory> [--detailed]")
		fmt.Println("       go run dvd_metadata.go <xml_file> [--detailed]")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	detailed := len(os.Args) > 2 && os.Args[2] == "--detailed"

	// Check if the argument is a directory or a file
	info, err := os.Stat(sourcePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var xmlFiles []string

	if info.IsDir() {
		// Process all XML files in the directory
		pattern := filepath.Join(sourcePath, "*.xml")
		xmlFiles, err = filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("Error finding XML files: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Process single file
		xmlFiles = []string{sourcePath}
	}

	if len(xmlFiles) == 0 {
		fmt.Printf("No XML files found in %s\n", sourcePath)
		os.Exit(1)
	}

	fmt.Printf("Found %d XML files to process\n", len(xmlFiles))

	for _, xmlFile := range xmlFiles {
		dvdData, err := dvd.ParseFile(xmlFile)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", xmlFile, err)
			continue
		}

		printDVDSummary(filepath.Base(xmlFile), dvdData)

		// If detailed mode is enabled, show detailed info for the longest track
		if detailed {
			longestTrack := dvdData.GetLongestTrack()
			if longestTrack != nil {
				printDetailedTrackInfo(*longestTrack)
			}
		}
	}
}
