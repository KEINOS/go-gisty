package gisty

import (
	"time"

	"github.com/KEINOS/go-gisty/gisty/internal/gistinfo"
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
	info, err := gistinfo.Parse(line)
	if err != nil {
		return GistInfo{}, WrapIfErr(err, "failed to parse line as GistInfo")
	}

	return GistInfo{
		GistID:      info.GistID,
		Description: info.Description,
		Files:       info.Files,
		IsPublic:    info.IsPublic,
		UpdatedAt:   info.UpdatedAt,
	}, nil
}
