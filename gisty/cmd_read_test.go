package gisty

import (
	"net/http"
	"testing"
	"time"

	"github.com/cli/cli/v2/pkg/cmd/gist/shared"
	"github.com/cli/cli/v2/pkg/cmd/gist/view"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGisty_Read_golden(t *testing.T) {
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
		}

		return gist, nil
	}

	// Instanciate Gisty object and execute the Read method.
	obj := NewGisty()
	gist, err := obj.Read("https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1")

	require.NoError(t, err)
	assert.Equal(t, "5b10b34f87955dfc86d310cd623a61d1", gist.ID)
	assert.Equal(t, "this is a dummy gist", gist.Description)
	assert.Equal(t, "https://gist.github.com/5b10b34f87955dfc86d310cd623a61d1", gist.HTMLURL)
	assert.Equal(t, 1, len(gist.Files))
	assert.Equal(t, "file1.txt", gist.Files["file1.txt"].Filename)
	assert.Equal(t, "text/plain", gist.Files["file1.txt"].Type)
	assert.Equal(t, "Text", gist.Files["file1.txt"].Language)
}

func TestGisty_Read_on_error(t *testing.T) {
	obj := NewGisty()

	obj.AltFunctions.Read = func(*view.ViewOptions) error {
		return errors.New("forced error")
	}

	gist, err := obj.Read("5b10b34f87955dfc86d310cd623a61d1")

	require.Error(t, err)
	require.Nil(t, gist, "gist should be nil on error")
	require.Contains(t, err.Error(), "failed to read gist")
	require.Contains(t, err.Error(), "forced error")
}

func TestGisty_Read_invalid_url(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(client *http.Client, hostname string, gistID string) (*shared.Gist, error) {
		return nil, nil
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

func TestGisty_Read_empty_gist(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(client *http.Client, hostname string, gistID string) (*shared.Gist, error) {
		return nil, nil
	}

	obj := NewGisty()

	// URL with control characters.
	gist, err := obj.Read("")

	require.Error(t, err)
	require.Nil(t, gist, "gist should be nil on error")
	require.Contains(t, err.Error(), "failed to read gist")
	require.Contains(t, err.Error(), "failed to execute readRun function")
	require.Contains(t, err.Error(), "no gist specified")
}

// ----------------------------------------------------------------------------
//  readRun()
// ----------------------------------------------------------------------------

func Test_readRun_fail_create_http_client(t *testing.T) {
	opts := &view.ViewOptions{
		Selector: "5b10b34f87955dfc86d310cd623a61d1",
		HttpClient: func() (*http.Client, error) {
			return nil, errors.New("forced error")
		},
	}

	gist, err := readRun(opts)

	require.Error(t, err)
	require.Nil(t, gist, "gist should be nil on error")
	require.Contains(t, err.Error(), "failed to create http client")
	require.Contains(t, err.Error(), "forced error")
}

func Test_readRun_fail_create_config(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	oldForceFailReadConf := forceFailReadConf
	defer func() {
		sharedGetGist = oldSharedGetGist
		forceFailReadConf = oldForceFailReadConf
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(client *http.Client, hostname string, gistID string) (*shared.Gist, error) {
		return nil, nil
	}

	forceFailReadConf = true

	obj := NewGisty()

	gist, err := obj.Read("5b10b34f87955dfc86d310cd623a61d1")

	require.Error(t, err)
	require.Nil(t, gist, "gist should be nil on error")
	require.Contains(t, err.Error(), "failed to read option config")
	require.Contains(t, err.Error(), "forced error")
}

func Test_readRun_fail_get_gist(t *testing.T) {
	oldSharedGetGist := sharedGetGist
	defer func() {
		sharedGetGist = oldSharedGetGist
	}()

	// Mock the shared.GetGist(sharedGetGist) function and return a dummy gist.
	sharedGetGist = func(client *http.Client, hostname string, gistID string) (*shared.Gist, error) {
		return nil, errors.New("forced error")
	}

	obj := NewGisty()

	gist, err := obj.Read("5b10b34f87955dfc86d310cd623a61d1")

	require.Error(t, err)
	require.Nil(t, gist, "gist should be nil on error")
	require.Contains(t, err.Error(), "failed to get gist")
	require.Contains(t, err.Error(), "forced error")
}
