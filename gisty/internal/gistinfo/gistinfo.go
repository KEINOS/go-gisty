// Package gistinfo parses the list output produced by gh gist list.
package gistinfo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	errEmptyLine         = errors.New("empty line")
	errMissingChunks     = errors.New("missing number of chunks")
	errInvalidVisibility = errors.New("failed to parse isPublic")
)

// Info holds parsed information about a gist.
type Info struct {
	UpdatedAt   time.Time
	GistID      string
	Description string
	Files       int
	IsPublic    bool
}

// Parse creates an Info value from one tab-separated gh gist list line.
func Parse(line string) (Info, error) {
	if line == "" {
		return Info{}, errEmptyLine
	}

	chunks := strings.Split(line, "\t")

	const chunkCount = 5

	if len(chunks) < chunkCount {
		return Info{}, fmt.Errorf(
			"%w: %d\ngiven line: %#v",
			errMissingChunks,
			len(chunks),
			line,
		)
	}

	files, err := ExtractFileNum(chunks[2])
	if err != nil {
		return Info{}, fmt.Errorf("failed to extract number of files: %w", err)
	}

	isPublic, err := ParseIsPublic(chunks[3])
	if err != nil {
		return Info{}, fmt.Errorf("failed to extract visibility: %w", err)
	}

	updatedAt, err := ParseTime(chunks[4])
	if err != nil {
		return Info{}, fmt.Errorf("failed to parse updatedAt: %w", err)
	}

	return Info{
		GistID:      chunks[0],
		Description: chunks[1],
		Files:       files,
		IsPublic:    isPublic,
		UpdatedAt:   updatedAt,
	}, nil
}

// ExtractFileNum parses a field formatted as "1 file" or "2 files".
func ExtractFileNum(chunk string) (int, error) {
	trimmed := strings.TrimSuffix(chunk, "files")
	trimmed = strings.TrimSuffix(trimmed, "file")
	trimmed = strings.TrimSpace(trimmed)

	numFiles, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, fmt.Errorf("failed to parse number of files from: %#v: %w", chunk, err)
	}

	return numFiles, nil
}

// ParseIsPublic parses the visibility field.
func ParseIsPublic(chunk string) (bool, error) {
	switch chunk {
	case "public":
		return true, nil
	case "secret":
		return false, nil
	default:
		return false, fmt.Errorf("%w from: %q", errInvalidVisibility, chunk)
	}
}

// ParseTime parses an RFC3339 timestamp.
func ParseTime(chunk string) (time.Time, error) {
	updatedAt, err := time.Parse(time.RFC3339, chunk)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time from: %q: %w", chunk, err)
	}

	return updatedAt, nil
}
