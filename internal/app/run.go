package app

import (
	"fmt"
	"io"

	"github.com/containeroo/tinyflags"
	"github.com/gi8lino/unveil/internal/collector"
	"github.com/gi8lino/unveil/internal/extract"
	"github.com/gi8lino/unveil/internal/flag"
	"github.com/gi8lino/unveil/internal/output"
)

// Run parses flags, builds specs, resolves values, and writes KEY=VAL lines.
func Run(
	version, commit string,
	args []string,
	w io.Writer,
) error {
	// parse flags (no side-effects here)
	flags, err := flag.ParseFlags(args, version, commit)
	if err != nil {
		if tinyflags.IsHelpRequested(err) || tinyflags.IsVersionRequested(err) {
			_, _ = fmt.Fprintln(w, err)
			return nil
		}
		return err
	}

	// collect specs
	specs, err := collector.Collect(&flags)
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return nil
	}

	// resolve values
	kv, err := extract.ExtractAll(specs)
	if err != nil {
		return err
	}

	// write output (atomic file or stdout)
	if flags.Output != "" {
		return output.WriteEnvLinesAtomic(flags.Output, kv, flags.Export)
	}

	return output.WriteEnvLines(w, kv, flags.Export)
}
