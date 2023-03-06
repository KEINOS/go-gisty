package gisty

import (
	"strconv"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
//  Type: GistInfo
// ----------------------------------------------------------------------------

// GistInfo holds information about a gist.
type GistInfo struct {
	UpdatedAt   time.Time // UpdatedAt is the time when the gist was last updated.
	GistID      string    // GistID is the ID of the gist.
	Description string    // Description is the description of the gist.
	Files       int       // Files is the number of files in the gist.
	IsPublic    bool      // IsPublic is true if the gist is public.
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// NewGistInfo creates a new GistInfo instance from the given input.
//
// The input line must be tab-separated and have the following format:
//
//		[gistID]\t[description]\t[n] <file|files>\t<public|secret>\t[updatedAt]
//	 e.g.
//		`1234567890abcdef	my gist	2 files	public	2018-01-01T00:00:00Z`
func NewGistInfo(line string) (GistInfo, error) {
	if line == "" {
		return GistInfo{}, NewErr("empty line")
	}

	chunks := strings.Split(line, "\t")
	numChumkMax := 5

	if len(chunks) < numChumkMax {
		return GistInfo{}, NewErr(
			"missing number of chunks: %d\ngiven line: %#v",
			len(chunks),
			line,
		)
	}

	files, err := extractFileNum(chunks[2])
	if err != nil {
		return GistInfo{}, WrapIfErr(err, "failed to extract number of files")
	}

	isPublic, err := parseIsPublic(chunks[3])
	if err != nil {
		return GistInfo{}, WrapIfErr(err, "failed to extract visibility")
	}

	updatedAt, err := parseTime(chunks[4])
	if err != nil {
		return GistInfo{}, WrapIfErr(err, "failed to parse updatedAt")
	}

	return GistInfo{
		GistID:      chunks[0],
		Description: chunks[1],
		Files:       files,
		IsPublic:    isPublic,
		UpdatedAt:   updatedAt,
	}, nil
}

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------

// extractFileNum extracts the file number from chunk.
//
//	e.g.) "1 file" -> 1
func extractFileNum(chunk string) (int, error) {
	trimmed := strings.TrimSuffix(chunk, "files")
	trimmed = strings.TrimSuffix(trimmed, "file")
	trimmed = strings.TrimSpace(trimmed)

	numFiles, err := strconv.Atoi(trimmed)

	return numFiles, WrapIfErr(err, "failed to parse number of files from: %#v", chunk)
}

// parseIsPublic parses the given string and returns true if the string is "public".
// If the string is "secret", it returns false. Otherwise, it returns an error.
func parseIsPublic(chunk string) (bool, error) {
	switch chunk {
	case "public":
		return true, nil
	case "secret":
		return false, nil
	default:
		return false, NewErr("failed to parse isPublic from: " + chunk)
	}
}

// parseTime parses the given RFC3339 time string and returns a time.Time instance.
func parseTime(chunk string) (time.Time, error) {
	// timeLayout := "2006-01-02T15:04:05Z"
	timeLayout := time.RFC3339

	updatedAt, err := time.Parse(timeLayout, chunk)
	if err != nil {
		return time.Time{}, WrapIfErr(err, "failed to parse time from: "+chunk)
	}

	return updatedAt, nil
}
