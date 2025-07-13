package gisty

import (
	"net/http"
	"testing"
	"time"

	"github.com/cli/cli/v2/pkg/cmd/gist/shared"
	"github.com/cli/cli/v2/pkg/cmd/gist/view"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func TestGisty_Read_golden(t *testing.T) {
	oldSharedGetGist := sharedGetGist

	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(_ *http.Client, _ string, gistID string) (*shared.Gist, error) {
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

	// Instantiate Gisty object and execute the Read method.
	obj := NewGisty()
	gist, err := obj.Read("https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1")

	require.NoError(t, err)
	assert.Len(t, gist.Files, 1)
	assert.Equal(t, "5b10b34f87955dfc86d310cd623a61d1", gist.ID)
	assert.Equal(t, "this is a dummy gist", gist.Description)
	assert.Equal(t, "https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1", gist.HTMLURL)
	assert.Equal(t, "file1.txt", gist.Files["file1.txt"].Filename)
	assert.Equal(t, "text/plain", gist.Files["file1.txt"].Type)
	assert.Equal(t, "Text", gist.Files["file1.txt"].Language)
}

func TestGisty_Read_on_error(t *testing.T) {
	t.Parallel()

	obj := NewGisty()

	obj.AltFunctions.Read = func(*view.ViewOptions) error {
		return NewErr("forced error")
	}

	gist, err := obj.Read("5b10b34f87955dfc86d310cd623a61d1")

	require.Error(t, err)
	require.Nil(t, gist, "gist should be nil on error")
	require.Contains(t, err.Error(), "failed to read gist")
	require.Contains(t, err.Error(), "forced error")
}

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func TestGisty_Read_invalid_url(t *testing.T) {
	oldSharedGetGist := sharedGetGist

	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(_ *http.Client, _ string, _ string) (*shared.Gist, error) {
		return new(shared.Gist), nil
	}

	obj := NewGisty()

	// URL with control characters.
	gist, err := obj.Read("https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1\n")

	require.Error(t, err)
	require.Nil(t, gist, "gist should be nil on error")
	require.Contains(t, err.Error(), "failed to read gist")
	require.Contains(t, err.Error(), "failed to execute readRun function")
	require.Contains(t, err.Error(), "failed to parse gist ID from URL")
}

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func TestGisty_Read_empty_gist(t *testing.T) {
	oldSharedGetGist := sharedGetGist

	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(_ *http.Client, _ string, _ string) (*shared.Gist, error) {
		return new(shared.Gist), nil
	}

	obj := NewGisty()

	// Empty gist URL
	gist, err := obj.Read("")

	require.Error(t, err)
	require.Nil(t, gist, "returned gist object should be nil on error")
	require.Contains(t, err.Error(), "failed to read gist")
	require.Contains(t, err.Error(), "failed to execute readRun function")
	require.Contains(t, err.Error(), "no gist specified")
}

// ----------------------------------------------------------------------------
//  readRun()
// ----------------------------------------------------------------------------

func Test_readRun_fail_create_http_client(t *testing.T) {
	t.Parallel()

	//nolint:exhaustruct // Missing fields are ok here.
	opts := &view.ViewOptions{
		Selector: "5b10b34f87955dfc86d310cd623a61d1",
		HttpClient: func() (*http.Client, error) {
			return nil, NewErr("forced error")
		},
	}

	gist, err := readRun(opts)

	require.Error(t, err)
	require.Nil(t, gist, "returned gist object should be nil on error")
	require.Contains(t, err.Error(), "failed to create http client")
	require.Contains(t, err.Error(), "forced error")
}

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func Test_readRun_fail_create_config(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	oldForceFailReadConf := forceFailReadConf

	defer func() {
		sharedGetGist = oldSharedGetGist
		forceFailReadConf = oldForceFailReadConf
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(_ *http.Client, _ string, _ string) (*shared.Gist, error) {
		return new(shared.Gist), nil
	}

	forceFailReadConf = true

	obj := NewGisty()

	gist, err := obj.Read("5b10b34f87955dfc86d310cd623a61d1")

	require.Error(t, err)
	require.Nil(t, gist, "returned gist object should be nil on error")
	require.Contains(t, err.Error(), "failed to read option config")
	require.Contains(t, err.Error(), "forced error")
}

//nolint:paralleltest // Do not parallelize due to mocking global function variables.
func Test_readRun_fail_get_gist(t *testing.T) {
	oldSharedGetGist := sharedGetGist

	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(_ *http.Client, _ string, _ string) (*shared.Gist, error) {
		return nil, NewErr("forced error")
	}

	obj := NewGisty()

	gist, err := obj.Read("5b10b34f87955dfc86d310cd623a61d1")

	require.Error(t, err)
	require.Nil(t, gist, "returned gist object should be nil on error")
	require.Contains(t, err.Error(), "failed to get gist")
	require.Contains(t, err.Error(), "forced error")
}
