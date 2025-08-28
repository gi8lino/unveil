package output

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errWriter struct {
	err error
}

func (e *errWriter) Write(p []byte) (int, error) { return 0, e.err }

func TestWriteEnvLines(t *testing.T) {
	t.Parallel()

	t.Run("Empty map writes nothing (no export)", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		err := WriteEnvLines(&buf, map[string]string{}, false)
		require.NoError(t, err)
		assert.Equal(t, "", buf.String())
	})

	t.Run("Single key=value (no export)", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		err := WriteEnvLines(&buf, map[string]string{"A": "1"}, false)
		require.NoError(t, err)
		assert.Equal(t, "A=1\n", buf.String())
	})

	t.Run("Multiple keys sorted by key (no export)", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		in := map[string]string{
			"Z": "last",
			"A": "first",
			"M": "middle",
		}
		err := WriteEnvLines(&buf, in, false)
		require.NoError(t, err)
		want := "A=first\nM=middle\nZ=last\n"
		assert.Equal(t, want, buf.String())
	})

	t.Run("Values are written verbatim (no quoting/escaping) (no export)", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		in := map[string]string{
			"SPECIAL": `a b "c" $d \ e`,
		}
		err := WriteEnvLines(&buf, in, false)
		require.NoError(t, err)
		assert.Equal(t, `SPECIAL=a b "c" $d \ e`+"\n", buf.String())
	})

	t.Run("Writer error is propagated (no export)", func(t *testing.T) {
		t.Parallel()
		w := &errWriter{err: errors.New("sink is broken")}
		err := WriteEnvLines(w, map[string]string{"A": "1"}, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sink is broken")
	})

	t.Run("Export prefix applied", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		err := WriteEnvLines(&buf, map[string]string{"A": "1", "B": "2"}, true)
		require.NoError(t, err)
		// Sorted keys, each with "export " prefix
		assert.Equal(t, "export A=1\nexport B=2\n", buf.String())
	})
}

func TestWriteEnvLinesAtomic(t *testing.T) {
	t.Parallel()

	t.Run("Writes file atomically (no export)", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "out.env")

		err := WriteEnvLinesAtomic(path, map[string]string{"B": "2", "A": "1"}, false)
		require.NoError(t, err)

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "A=1\nB=2\n", string(got))
	})

	t.Run("Creates parent directories", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "nested", "deeper", "out.env")

		err := WriteEnvLinesAtomic(path, map[string]string{"K": "V"}, false)
		require.NoError(t, err)

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "K=V\n", string(got))
	})

	t.Run("Export prefix applied", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "export.env")

		err := WriteEnvLinesAtomic(path, map[string]string{"A": "1", "B": "2"}, true)
		require.NoError(t, err)

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "export A=1\nexport B=2\n", string(got))
	})

	t.Run("Overwrites existing file atomically", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "out.env")

		// Seed with old content
		require.NoError(t, os.WriteFile(path, []byte("OLD=1\n"), 0o666))

		// Write new content
		err := WriteEnvLinesAtomic(path, map[string]string{"NEW": "2"}, false)
		require.NoError(t, err)

		got, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "NEW=2\n", string(got))
	})
}
