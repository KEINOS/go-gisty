package gisty

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cli/cli/v2/pkg/cmd/api"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// For the example usge of Comments method, see: example_test.go

func TestGisty_Comments_bad_id(t *testing.T) {
	objGisty := NewGisty()
	listComments, err := objGisty.Comments("\n\t")

	require.Error(t, err,
		"invalid gist ID (non hexdecimal characters) should return error")
	require.Nil(t, listComments,
		"returned slice of comment objects should be nil on error")
}

func TestGisty_comments_golden(t *testing.T) {
	obj := NewGisty()

	// Dummy function to avoid calling the actual GitHub API during test/example.
	// Usually, you do not need to set this.
	obj.AltFunctions.Comments = func(*api.ApiOptions) error {
		resp, err := json.Marshal([]Comment{DummyComment})
		require.NoError(t, err,
			"failed to create dummy data during test setup")

		// Mock the GitHub API response.
		fmt.Fprint(
			obj.Stdout,
			string(resp),
		)

		return nil
	}

	listComments, err := obj.Comments("abcdef1234567890")

	require.NoError(t, err,
		"no error should be returned when the gist ID is valid")
	require.NotNil(t, listComments,
		"returned slice of comment objects should not be nil")
}

func TestGisty_comments_execute_error(t *testing.T) {
	obj := NewGisty()

	// Dummy function to avoid calling the actual GitHub API during test/example.
	// Usually, you do not need to set this.
	obj.AltFunctions.Comments = func(*api.ApiOptions) error {
		return errors.New("forced error")
	}

	listComments, err := obj.Comments("abcdef1234567890")

	require.Error(t, err,
		"invalid gist ID (non hexdecimal characters) should return error")
	require.Nil(t, listComments,
		"returned slice of comment objects should be nil on error")
	require.Contains(t, err.Error(), "failed to execute GitHub API request",
		"error message should contain the error reason")
	require.Contains(t, err.Error(), "forced error",
		"error message should contain the underlying error message")
}

func TestGisty_comments_malformed_response(t *testing.T) {
	obj := NewGisty()

	// Dummy function to avoid calling the actual GitHub API during test/example.
	// Usually, you do not need to set this.
	obj.AltFunctions.Comments = func(*api.ApiOptions) error {
		// Mock the GitHub API response.
		fmt.Fprint(
			obj.Stdout,
			"this is an invalid JSON: which, is, not, valid, JSON, at, all",
		)

		return nil
	}

	listComments, err := obj.Comments("abcdef1234567890")

	require.Error(t, err,
		"invalid gist ID (non hexdecimal characters) should return error")
	require.Nil(t, listComments,
		"returned slice of comment objects should be nil on error")
	require.Contains(t, err.Error(), "failed to parse GitHub API response. malformed JSON",
		"error message should contain the error reason")
}
