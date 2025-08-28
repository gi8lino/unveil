package output

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// WriteEnvLines prints KEY=VALUE lines, sorted by KEY.
// If export is true, each line is prefixed with "export ".
func WriteEnvLines(w io.Writer, kv map[string]string, export bool) error {
	prefix := ""
	if export {
		prefix = "export "
	}

	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if _, err := fmt.Fprintf(w, "%s%s=%s\n", prefix, k, kv[k]); err != nil {
			return err
		}
	}
	return nil
}

// WriteEnvLinesAtomic writes KEY=VALUE lines atomically to path.
// It creates parent directories, writes to a temp file, fsyncs, and renames.
func WriteEnvLinesAtomic(path string, kv map[string]string, export bool) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating output directory %q: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, ".unveil-*")
	if err != nil {
		return fmt.Errorf("creating temp file in %q: %w", dir, err)
	}
	tmpPath := tmp.Name()
	// ensure cleanup on any early return
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	if err := WriteEnvLines(tmp, kv, export); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("syncing temp file %q: %w", tmpPath, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file %q: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("renaming %q to %q: %w", tmpPath, path, err)
	}
	return nil
}
