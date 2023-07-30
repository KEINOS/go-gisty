package gisty

import (
	"net/http"
	"testing"
	"time"

	"github.com/cli/cli/v2/pkg/cmd/gist/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func TestGisty_ReadFile_golden(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	expectContent := "This is the content of the file1."

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(client *http.Client, hostname string, gistID string) (*shared.Gist, error) {
		// Dummy files in a gist.
		files := map[string]*shared.GistFile{
			"file1.txt": {
				Filename: "file1.txt",
				Type:     "text/plain",
				Language: "Text",
				Content:  expectContent,
			},
		}

		// Create a dummy Gist object.
		gist := &shared.Gist{
			ID:          gistID,
			Description: "this is a dummy gist",
			Files:       files,
			UpdatedAt:   time.Now(),
			Public:      true,
			HTMLURL:     "https://gist.github.com/" + gistID,
			Owner:       nil,
		}

		return gist, nil
	}

	actualContent, err := NewGisty().ReadFile(
		"https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1",
		"file1.txt",
	)

	require.NoError(t, err)
	assert.Equal(t, expectContent, string(actualContent))
}

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func TestGisty_ReadFile_read_fail(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(client *http.Client, hostname string, gistID string) (*shared.Gist, error) {
		return nil, NewErr("forced error")
	}

	data, err := NewGisty().ReadFile(
		"https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1",
		"file1.txt",
	)

	require.Error(t, err,
		"it should error out when reading the gist info fails")
	require.Nil(t, data,
		"returned data should be nil on error")
	require.Contains(t, err.Error(), "failed to read gist info",
		"it should contain the error reason")
	require.Contains(t, err.Error(), "forced error",
		"it should contain the underlying error reason")
}

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func TestGisty_ReadFile_file_not_found(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(client *http.Client, hostname string, gistID string) (*shared.Gist, error) {
		// Dummy files in a gist.
		files := map[string]*shared.GistFile{
			"file1.txt": {
				Filename: "file1.txt",
				Type:     "text/plain",
				Language: "Text",
				Content:  "This is the content of the file1.",
			},
		}

		// Create a dummy Gist object.
		gist := &shared.Gist{
			ID:          gistID,
			Description: "this is a dummy gist",
			Files:       files,
			UpdatedAt:   time.Now(),
			Public:      true,
			HTMLURL:     "https://gist.github.com/" + gistID,
			Owner:       nil,
		}

		return gist, nil
	}

	data, err := NewGisty().ReadFile(
		"https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1",
		"file2.txt", // non-existent file
	)

	require.Error(t, err,
		"it should error out when reading the gist info fails")
	require.Nil(t, data,
		"returned data should be nil on error")
	require.Contains(t, err.Error(), "file not found",
		"it should contain the error reason")
}
