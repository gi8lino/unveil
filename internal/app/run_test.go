package app

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errWriter struct{ err error }

func (e *errWriter) Write(p []byte) (int, error) { return 0, e.err }

func TestRun_HelpRequested(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	err := Run("1.0.0", "abc", []string{"--help"}, &out)
	require.NoError(t, err)
	// tinyflags prints a help message to stdout; just ensure something was printed
	assert.NotEmpty(t, out.String())
}

func TestRun_VersionRequested(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	err := Run("9.9.9", "deadbeef", []string{"--version"}, &out)
	require.NoError(t, err)
	// Version text should be printed; don't depend on exact wording
	assert.Equal(t, "9.9.9\n", out.String())
	assert.NotEmpty(t, out.String())
}

func TestRun_ParseError_MissingRequired(t *testing.T) {
	t.Parallel()

	t.Run("Missing required --json.id.select", func(t *testing.T) {
		var out bytes.Buffer
		// Missing required --json.id.select
		err := Run("v", "c", []string{"--json.id.path=./cfg.json"}, &out)
		require.Error(t, err)
		assert.Empty(t, out.String())
	})

	t.Run("No specs, no output", func(t *testing.T) {
		t.Parallel()

		var out bytes.Buffer
		err := Run("v", "c", []string{}, &out)
		require.NoError(t, err)
		assert.Equal(t, "", out.String())
	})

	t.Run("Successful multi-spec extraction with sorted output", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		// Prepare real files for resolvers
		// JSON
		jpath := filepath.Join(dir, "cfg.json")
		require.NoError(t, os.WriteFile(jpath, []byte(`{"server":{"host":"localhost"}}`), 0o666))

		// INI
		ipath := filepath.Join(dir, "conf.ini")
		require.NoError(t, os.WriteFile(ipath, []byte("[DB]\nUser=alice\n"), 0o666))

		// FILE (key=value)
		fpath := filepath.Join(dir, "app.env")
		require.NoError(t, os.WriteFile(fpath, []byte("TOKEN=secret\n"), 0o666))

		args := []string{
			"--json.web.path=" + jpath,
			"--json.web.select=server.host",
			"--json.web.as=HOST",

			"--ini.db.path=" + ipath,
			"--ini.db.select=DB.User",
			"--ini.db.as=DBUSER",

			"--file.env.path=" + fpath,
			"--file.env.select=TOKEN",
			"--file.env.as=APPTOKEN",
		}

		var out bytes.Buffer
		err := Run("v", "c", args, &out)
		require.NoError(t, err)

		// Output must be sorted by KEY alphabetically
		got := out.String()
		want := "" +
			"APPTOKEN=secret\n" +
			"DBUSER=alice\n" +
			"HOST=localhost\n"
		assert.Equal(t, want, got)
	})

	t.Run("Extraction error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		tpath := filepath.Join(dir, "cfg.toml")
		require.NoError(t, os.WriteFile(tpath, []byte(`[server] host = "h"`), 0o666))

		args := []string{
			"--toml.svc.path=" + tpath,
			"--toml.svc.select=server.port", // missing
			"--toml.svc.as=PORT",
		}

		var out bytes.Buffer
		err := Run("v", "c", args, &out)
		require.Error(t, err)
		assert.Empty(t, out.String())
	})

	t.Run("Output write error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		jpath := filepath.Join(dir, "cfg.json")
		require.NoError(t, os.WriteFile(jpath, []byte(`{"k":"v"}`), 0o666))

		args := []string{
			"--json.a.path=" + jpath,
			"--json.a.select=k",
			"--json.a.as=K",
		}

		w := &errWriter{err: errors.New("sink broken")}
		err := Run("v", "c", args, w)
		require.Error(t, err)
	})

	t.Run("Unknown group", func(t *testing.T) {
		t.Parallel()
		args := []string{
			"--unknown.a.path=unknown.txt",
			"--unknown.a.select=k",
			"--unknown.a.as=K",
		}

		var out bytes.Buffer
		err := Run("v", "c", args, &out)
		require.Error(t, err)
		assert.EqualError(t, err, "unknown dynamic group \"unknown\" in flag --unknown.a.path=unknown.txt\nunknown dynamic group \"unknown\" in flag --unknown.a.select=k\nunknown dynamic group \"unknown\" in flag --unknown.a.as=K")
	})

	t.Run("Write output to file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		src := filepath.Join(dir, "app.env")
		dst := filepath.Join(dir, "out.env")
		require.NoError(t, os.WriteFile(src, []byte("TOKEN=secret\n"), 0o666))

		args := []string{
			"--file.env.path=" + src,
			"--file.env.select=TOKEN",
			"--file.env.as=APPTOKEN",
			"--output=" + dst,
		}

		var out bytes.Buffer
		err := Run("v", "c", args, &out)
		require.NoError(t, err)

		got, err := os.ReadFile(dst)
		require.NoError(t, err)

		want := "APPTOKEN=secret\n"
		assert.Equal(t, want, string(got))
		assert.Equal(t, "", out.String())
	})
}
