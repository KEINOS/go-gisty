package gisty

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // This test replaces the package-level command executor.
func TestGisty_default_commands(t *testing.T) {
	stubGHCommand(t, false)

	obj := NewGisty()
	require.NoError(t, obj.Clone([]string{"dummy"}))

	obj = NewGisty()
	gistURL, err := obj.Create(CreateArgs{
		Description: "",
		FilePaths:   []string{"testdata/foo.md"},
		AsPublic:    false,
	})
	require.NoError(t, err)
	require.Equal(t, "https://gist.github.com/dummy", gistURL.String())

	obj = NewGisty()
	require.NoError(t, obj.Delete("dummy"))

	obj = NewGisty()
	gists, err := obj.List(ListArgs{
		Limit:      1,
		OnlyPublic: false,
		OnlySecret: false,
	})
	require.NoError(t, err)
	require.Len(t, gists, 1)

	obj = NewGisty()
	comments, err := obj.Comments("dummy")
	require.NoError(t, err)
	require.Empty(t, comments)

	obj = NewGisty()
	count, err := obj.Stargazer("dummy")
	require.NoError(t, err)
	require.Equal(t, 42, count)

	obj = NewGisty()
	_, err = obj.Update(NewUpdateArgs(t.TempDir()))
	require.NoError(t, err)
}

//nolint:paralleltest // This test replaces the package-level command executor.
func TestGisty_runGH_error(t *testing.T) {
	stubGHCommand(t, true)

	obj := NewGisty()

	for _, test := range []struct {
		name string
		run  func() error
	}{
		{name: "runGH", run: func() error { return obj.runGH("version") }},
		{name: "clone", run: func() error { return obj.Clone([]string{"dummy"}) }},
		{name: "create", run: func() error {
			_, err := obj.Create(CreateArgs{Description: "", FilePaths: nil, AsPublic: false})

			return err
		}},
		{name: "delete", run: func() error { return obj.Delete("dummy") }},
		{name: "list", run: func() error {
			_, err := obj.List(ListArgs{Limit: 1, OnlyPublic: false, OnlySecret: false})

			return err
		}},
		{name: "comments", run: func() error {
			_, err := obj.Comments("dummy")

			return err
		}},
		{name: "stargazer", run: func() error {
			_, err := obj.Stargazer("dummy")

			return err
		}},
		{name: "update", run: func() error {
			_, err := obj.Update(NewUpdateArgs(t.TempDir()))

			return err
		}},
	} {
		err := test.run()
		require.Error(t, err, test.name)
		require.Contains(t, err.Error(), "failed to execute gh command", test.name)
	}
}

//nolint:paralleltest // This test is executed as a subprocess by stubGHCommand.
func TestGHHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_GH_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	indexSeparator := slices.Index(args, "--")
	require.NotEqual(t, -1, indexSeparator)
	args = args[indexSeparator+1:]

	if slices.Contains(args, "forced-error") {
		t.Fatal("forced error")
	}

	switch strings.Join(args[:min(2, len(args))], " ") {
	case "gist create":
		_, err := fmt.Fprint(os.Stdout, "https://gist.github.com/dummy\n")
		require.NoError(t, err)
	case "gist list":
		_, err := fmt.Fprint(os.Stdout, "dummy\tdescription\t1 file\tpublic\t2026-05-31T00:00:00Z\n")
		require.NoError(t, err)
	case "api graphql":
		if slices.Contains(args, "--jq") {
			_, err := fmt.Fprint(os.Stdout, "[]\n")
			require.NoError(t, err)
		} else {
			_, err := fmt.Fprint(os.Stdout, "'42'")
			require.NoError(t, err)
		}
	case "repo sync":
		_, err := fmt.Fprint(os.Stdout, "✓ Synced\n")
		require.NoError(t, err)
	}

	os.Exit(0)
}

func stubGHCommand(t *testing.T, forceError bool) {
	t.Helper()

	oldExecCommandContext := execCommandContext

	t.Cleanup(func() {
		execCommandContext = oldExecCommandContext
	})

	execCommandContext = func(ctx context.Context, _ string, args ...string) *exec.Cmd {
		if forceError {
			args = []string{"forced-error"}
		}

		helperArgs := append([]string{"-test.run=TestGHHelperProcess", "--"}, args...)
		//nolint:gosec // The helper executes the current test binary with controlled arguments.
		cmd := exec.CommandContext(ctx, os.Args[0], helperArgs...)

		cmd.Env = append(os.Environ(), "GO_WANT_GH_HELPER_PROCESS=1")

		return cmd
	}
}
