package gisty

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Type: GistInfo
// ----------------------------------------------------------------------------

// GistInfo holds information about a gist.
type GistInfo struct {
	UpdatedAt   time.Time
	GistID      string
	Description string
	Files       int
	IsPublic    bool
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
		return GistInfo{}, errors.New("empty line")
	}

	chunks := strings.Split(line, "\t")
	numChumkMax := 5

	if len(chunks) < numChumkMax {
		return GistInfo{}, errors.Errorf(
			"missing number of chunks: %d\ngiven line: %#v",
			len(chunks),
			line,
		)
	}

	files, err := extractFileNum(chunks[2])
	if err != nil {
		return GistInfo{}, err
	}

	isPublic, err := parseIsPublic(chunks[3])
	if err != nil {
		return GistInfo{}, err
	}

	updatedAt, err := parseTime(chunks[4])
	if err != nil {
		return GistInfo{}, err
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
	chunk = strings.TrimSuffix(chunk, "files")
	chunk = strings.TrimSuffix(chunk, "file")
	chunk = strings.TrimSpace(chunk)

	numFiles, err := strconv.Atoi(chunk)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse number of files from: "+chunk)
	}

	return numFiles, nil
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
		return false, errors.New("failed to parse isPublic from: " + chunk)
	}
}

// parseTime parses the given RFC3339 time string and returns a time.Time instance.
func parseTime(chunk string) (time.Time, error) {
	// timeLayout := "2006-01-02T15:04:05Z"
	timeLayout := time.RFC3339

	updatedAt, err := time.Parse(timeLayout, chunk)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to parse time from: "+chunk)
	}

	return updatedAt, nil
}
