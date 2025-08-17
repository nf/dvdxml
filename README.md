# DVD Metadata Parser

A Go library and command-line program that parses XML files containing DVD disc metadata and extracts comprehensive information about tracks, chapters, audio streams, and subtitle streams.

## Library Usage

The DVD parsing functionality is available as a reusable Go package:

```go
import "dvd-metadata-parser/dvd"

// Parse a DVD metadata file
dvdData, err := dvd.ParseFile("path/to/file.xml")
if err != nil {
    log.Fatal(err)
}

// Access DVD information
fmt.Printf("Title: %s\n", dvdData.Title)
fmt.Printf("Tracks: %d\n", len(dvdData.Tracks))
fmt.Printf("Total duration: %.2f minutes\n", dvdData.GetTotalDuration()/60)

// Get the longest track
longestTrack := dvdData.GetLongestTrack()
if longestTrack != nil {
    fmt.Printf("Longest track: #%d (%.2f minutes)\n", 
        longestTrack.Index, longestTrack.Length/60)
}

// Get available languages
audioLangs := dvdData.GetAudioLanguages()
subLangs := dvdData.GetSubtitleLanguages()
```

See [example_usage.go.example](example_usage.go.example) for a complete example.

## Features

- **Complete DVD Structure Parsing**: Extracts all metadata from lsdvd XML output
- **Track Information**: Video format, resolution, aspect ratio, frame rate, and duration
- **Audio Stream Details**: Language, format, frequency, channel count, and stream IDs
- **Subtitle Stream Information**: Language codes, content type, and stream IDs
- **Chapter Breakdown**: Individual chapter durations and cell information
- **Error Handling**: Robust parsing with automatic correction of malformed XML entities
- **Flexible Input**: Process single files or entire directories
- **Detailed Mode**: Extended information display for the longest track

## Usage

### Basic usage
```bash
go run dvd_metadata.go <xml_file>
go run dvd_metadata.go source/s1d1.xml

go run dvd_metadata.go <directory>
go run dvd_metadata.go source
```

### Show detailed information for the longest track
```bash
go run dvd_metadata.go -detailed source/s1d1.xml
go run dvd_metadata.go -detailed source
```

### Find episodes of specific duration
```bash
# Find content around 40 minutes (±5 minutes by default)
go run dvd_metadata.go -episodes 40 source

# Find content around 22 minutes with custom tolerance
go run dvd_metadata.go -episodes 22 -tolerance 3 source

# Find movie-length content
go run dvd_metadata.go -episodes 90 -tolerance 15 source
```

### View help
```bash
go run dvd_metadata.go -help
```

## Output Format

The program provides a structured summary for each DVD:

```
=== s1d1.xml ===
Device: ./s1d1/Law And Order Svu
Title: unknown
Provider ID: 
Number of tracks: 10
Longest track: 5

  Track 1:
    Length: 2500.56 seconds (41.68 minutes)
    Resolution: 720x576
    Aspect: 4/3
    Format: PAL @ 25.00 fps
    Chapters: 5
    Audio streams: 2
    Subtitle streams: 4
      Audio 1: English (en) - ac3, 48000 Hz, 2 channels
      Audio 2: Francais (fr) - ac3, 48000 Hz, 2 channels
      Subtitle 1: English (en)
      Subtitle 2: Francais (fr)
```

### Detailed Mode Output

With `--detailed` flag, the program shows comprehensive information about the longest track:

```
--- Detailed Track 5 Information ---
Length: 9857.28 seconds
Video: PAL, 720x576, 4/3, 25.00 fps
VTS: 1, TTN: 5

Audio Streams (2):
  [1] English (en): ac3, 48000 Hz, 2 ch, Stream ID: 0x80
  [2] Francais (fr): ac3, 48000 Hz, 2 ch, Stream ID: 0x81

Subtitle Streams (4):
  [1] English (en): Undefined, Stream ID: 0x20
  [2] Francais (fr): Undefined, Stream ID: 0x21
  [3] Nederlands (nl): Undefined, Stream ID: 0x22
  [4] Francais (fr): Undefined, Stream ID: 0x23

Chapters (17):
  Chapter 1: 735.20 seconds (starts at cell 1)
  Chapter 2: 423.20 seconds (starts at cell 2)
  ...
```

### Episodes Mode Output

With `-episodes` flag, the program finds tracks and chapters around the specified duration:

```
=== s1d1.xml - ~40 Minute Content ===
Looking for content between 35.0-45.0 minutes...

  ✓ Track 1: 41.68 minutes (2500.56 seconds)
    Resolution: 720x576, Format: PAL @ 25.00 fps
    Audio: 2 streams, Subtitles: 4 streams, Chapters: 5

  ✓ Track 2: 41.00 minutes (2459.88 seconds)
    Resolution: 720x576, Format: PAL @ 25.00 fps
    Audio: 2 streams, Subtitles: 4 streams, Chapters: 5

Summary: 4 tracks and 0 chapters found around 40 minutes.
```

## Package API

The `dvd` package provides the following types and functions:

### Types
- **`DVD`**: Root container with device info and tracks
- **`Track`**: Individual title with video, audio, subtitle, and chapter data
- **`AudioStream`**: Audio track details (language, format, channels, etc.)
- **`SubtitleStream`**: Subtitle track information
- **`Chapter`**: Chapter timing and cell references
- **`Cell`**: DVD cell structure information
- **`Palette`**: Color palette data
- **`ContentMatch`**: Represents a track or chapter that matches duration criteria

### Functions
- **`ParseFile(filename string) (*DVD, error)`**: Parse DVD metadata from XML file
- **`ParseBytes(data []byte) (*DVD, error)`**: Parse DVD metadata from XML byte data

### Methods on DVD
- **`GetLongestTrack() *Track`**: Returns the longest track
- **`GetTrackByIndex(index int) *Track`**: Returns track by index (1-based)
- **`GetTotalDuration() float64`**: Returns total duration of all tracks
- **`GetAudioLanguages() []string`**: Returns unique audio languages
- **`GetSubtitleLanguages() []string`**: Returns unique subtitle languages
- **`FindFortyMinuteContent() []ContentMatch`**: Finds tracks/chapters around 40 minutes
- **`FindContentAroundDuration(targetMinutes, toleranceMinutes float64) []ContentMatch`**: Finds content around any duration

## XML Format Support

The program is designed to parse XML files generated by the `lsdvd` tool, which analyzes DVD structure. Key supported elements:

- `<lsdvd>` - Root element
- `<track>` - Individual DVD tracks/titles
- `<audio>` - Audio stream information
- `<subp>` - Subtitle stream data
- `<chapter>` - Chapter definitions
- `<cell>` - DVD cell information

## Error Handling

The program includes robust error handling:

- **Malformed XML Entities**: Automatically fixes common issues like `Pan&Scan` → `Pan&amp;Scan`
- **Missing Files**: Graceful error messages for non-existent files
- **Invalid XML**: Clear error reporting with file names and line numbers
- **Partial Failures**: Continues processing other files even if some fail

## Testing

Run the test suite:

```bash
go test -v
```

The tests validate:
- Individual XML file parsing
- All XML files in the source directory
- Specific field values and data integrity
- Error handling for malformed XML

## Dependencies

- Go 1.21 or later
- Standard library only (no external dependencies)

## File Structure

```
.
├── dvd/                 # DVD parsing package
│   ├── parser.go        # Core parsing logic and types
│   └── parser_test.go   # Package-specific tests
├── dvd_metadata.go      # Command-line program
├── dvd_metadata_test.go # Integration tests
├── example_usage.go.example # Library usage example
├── go.mod              # Go module definition
├── README.md           # This documentation
├── source/             # Directory containing XML files
│   ├── s1d1.xml
│   ├── s1d2.xml
│   └── ...
└── output_example.txt  # Example output
```

## Example XML Structure

The program expects XML files with structure similar to:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<lsdvd>
    <device>./s1d1/Law And Order Svu</device>
    <title>unknown</title>
    <vmg_id>DVDVIDEO-VMG</vmg_id>
    <provider_id></provider_id>
    <track>
        <ix>1</ix>
        <length>2500.560</length>
        <fps>25.00</fps>
        <format>PAL</format>
        <width>720</width>
        <height>576</height>
        <audio>
            <ix>1</ix>
            <language>English</language>
            <format>ac3</format>
            <frequency>48000</frequency>
            <channels>2</channels>
        </audio>
        <chapter>
            <ix>1</ix>
            <length>735.200</length>
            <startcell>1</startcell>
        </chapter>
    </track>
</lsdvd>
```
