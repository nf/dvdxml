package main

import (
	"dvd-metadata-parser/dvd"
	"flag"
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

// generateFFmpegCommand generates an FFmpeg command to extract a track or chapter
func generateFFmpegCommand(match dvd.ContentMatch, dvdPath, outputPrefix string) string {
	if match.Type == "track" {
		// Extract entire track using dvdvideo demuxer
		outputFile := fmt.Sprintf("%s_track_%02d.mkv", outputPrefix, match.Track.Index)
		// Use dvdvideo:path and specify the title (track) to extract
		return fmt.Sprintf("ffmpeg -f dvdvideo -i '%s' -title %d -map 0 -c copy %q",
			dvdPath, match.Track.Index, outputFile)
	} else {
		// Extract specific chapter range - this is more complex and would need chapter timing
		outputFile := fmt.Sprintf("%s_track_%02d_chapter_%02d.mkv",
			outputPrefix, match.Track.Index, match.Chapter.Index)
		return fmt.Sprintf("ffmpeg -f dvdvideo -i '%s' -title %d -chapter_start %d -chapter_end %d -map 0 -c copy %q",
			dvdPath, match.Track.Index, match.Chapter.Index, match.Chapter.Index+1, outputFile)
	}
}

// extractDVDPath tries to extract the DVD path from device string
func extractDVDPath(device string) string {
	// Remove common prefixes like "./" and get just the directory name
	if len(device) >= 2 && device[:2] == "./" {
		return device[2:]
	}
	return device
}

// findEpisodeContent finds tracks and chapters around a specified duration
func findEpisodeContent(filename string, dvdData *dvd.DVD, targetMinutes, toleranceMinutes float64) {
	fmt.Printf("\n=== %s - ~%.0f Minute Content ===\n", filename, targetMinutes)
	fmt.Printf("Looking for content between %.1f-%.1f minutes...\n",
		targetMinutes-toleranceMinutes, targetMinutes+toleranceMinutes)

	matches := dvdData.FindContentAroundDuration(targetMinutes, toleranceMinutes)

	if len(matches) == 0 {
		fmt.Printf("  No tracks or chapters found around %.0f minutes.\n", targetMinutes)
		return
	}

	tracksFound := 0
	chaptersFound := 0
	currentTrack := -1

	for _, match := range matches {
		if match.Type == "track" {
			tracksFound++
			fmt.Printf("\n  ✓ Track %d: %.2f minutes (%.2f seconds)\n",
				match.Track.Index, match.Duration/60, match.Duration)
			fmt.Printf("    Resolution: %dx%d, Format: %s @ %.2f fps\n",
				match.Track.Width, match.Track.Height, match.Track.Format, match.Track.FPS)
			fmt.Printf("    Audio: %d streams, Subtitles: %d streams, Chapters: %d\n",
				len(match.Track.AudioStreams), len(match.Track.SubtitleStreams), len(match.Track.Chapters))
		} else if match.Type == "chapter" {
			chaptersFound++
			if match.Track.Index != currentTrack {
				currentTrack = match.Track.Index
				fmt.Printf("\n  Track %d chapters:\n", match.Track.Index)
				fmt.Printf("    Track length: %.2f minutes, Resolution: %dx%d\n",
					match.Track.Length/60, match.Track.Width, match.Track.Height)
			}
			fmt.Printf("    ✓ Chapter %d: %.2f minutes (%.2f seconds)\n",
				match.Chapter.Index, match.Duration/60, match.Duration)
		}
	}

	fmt.Printf("\nSummary: %d tracks and %d chapters found around %.0f minutes.\n",
		tracksFound, chaptersFound, targetMinutes)
}

func main() {
	// Define command line flags
	var (
		detailed  = flag.Bool("detailed", false, "Show detailed info for longest track")
		episodes  = flag.Float64("episodes", 0, "Find tracks/chapters around specified duration in minutes (e.g., 40)")
		tolerance = flag.Float64("tolerance", 5.0, "Tolerance in minutes for episode duration matching (default: 5)")
		ffmpeg    = flag.Bool("ffmpeg", false, "Generate FFmpeg commands to extract episodes (use with -episodes)")
		showHelp  = flag.Bool("help", false, "Show this help message")
	) // Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <source_directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s [flags] <xml_file>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s source/s1d1.xml                    # Basic summary\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -detailed source                   # Show detailed longest track info\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -episodes 40 source                # Find ~40 minute episodes\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -episodes 22 -tolerance 3 source   # Find ~22 minute episodes (±3 min)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -episodes 40 -ffmpeg source        # Generate FFmpeg commands for extraction\n", os.Args[0])
	}

	// Parse command line flags
	flag.Parse()

	// Show help if requested
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Check for required source path argument
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: Please specify exactly one source directory or XML file\n\n")
		flag.Usage()
		os.Exit(1)
	}

	sourcePath := flag.Arg(0)

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

	// Only show processing message in non-FFmpeg mode
	if !(*episodes > 0 && *ffmpeg) {
		fmt.Printf("Found %d XML files to process\n", len(xmlFiles))
	}

	for _, xmlFile := range xmlFiles {
		dvdData, err := dvd.ParseFile(xmlFile)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", xmlFile, err)
			continue
		}

		if *episodes > 0 {
			if *ffmpeg {
				// FFmpeg mode: only output commands
				matches := dvdData.FindContentAroundDuration(*episodes, *tolerance)
				if len(matches) > 0 {
					dvdPath := extractDVDPath(dvdData.Device)
					outputPrefix := fmt.Sprintf("%s_episodes", filepath.Base(xmlFile)[:len(filepath.Base(xmlFile))-4])
					for _, match := range matches {
						if match.Type == "track" {
							cmd := generateFFmpegCommand(match, dvdPath, outputPrefix)
							fmt.Println(cmd)
						}
					}
				}
			} else {
				findEpisodeContent(filepath.Base(xmlFile), dvdData, *episodes, *tolerance)
			}
		} else {
			printDVDSummary(filepath.Base(xmlFile), dvdData)

			// If detailed mode is enabled, show detailed info for the longest track
			if *detailed {
				longestTrack := dvdData.GetLongestTrack()
				if longestTrack != nil {
					printDetailedTrackInfo(*longestTrack)
				}
			}
		}
	}
}
