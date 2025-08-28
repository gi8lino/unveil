package collector

import (
	"fmt"
	"strings"

	"github.com/gi8lino/unveil/internal/flag"
	"github.com/gi8lino/unveil/internal/quote"
	"github.com/gi8lino/unveil/internal/spec"

	"github.com/containeroo/tinyflags"
)

// Collect flattens all dynamic group instances into a slice of ExtractSpec.
func Collect(flags *flag.Flags) ([]spec.ExtractSpec, error) {
	var out []spec.ExtractSpec

	for _, g := range flags.FlagSet.DynamicGroups() {
		groupName := g.Name() // one of: "json", "yaml", "toml", "ini", "file"

		for _, id := range g.Instances() {
			// Read per-instance values
			path := tinyflags.GetOrDefaultDynamic[string](g, id, "path")
			key := tinyflags.GetOrDefaultDynamic[string](g, id, "select")
			varName := tinyflags.GetOrDefaultDynamic[string](g, id, "as")
			quoteStr := tinyflags.GetOrDefaultDynamic[string](g, id, "quote")

			// Effective quoting: instance override if present, else global default
			quoteKind := flags.Quote
			if quoteStr != "" {
				quoteKind = quote.QuoteKind(quoteStr)
			}

			if varName == "" {
				varName = strings.ToUpper(id)
			}

			// Append with appropriate Kind
			switch groupName {
			case "json":
				out = append(out, spec.ExtractSpec{Kind: spec.KindJSON, Path: path, Key: key, Var: varName, Quote: quoteKind})
			case "yaml":
				out = append(out, spec.ExtractSpec{Kind: spec.KindYAML, Path: path, Key: key, Var: varName, Quote: quoteKind})
			case "toml":
				out = append(out, spec.ExtractSpec{Kind: spec.KindTOML, Path: path, Key: key, Var: varName, Quote: quoteKind})
			case "ini":
				out = append(out, spec.ExtractSpec{Kind: spec.KindINI, Path: path, Key: key, Var: varName, Quote: quoteKind})
			case "file":
				out = append(out, spec.ExtractSpec{Kind: spec.KindFILE, Path: path, Key: key, Var: varName, Quote: quoteKind})
			default:
				// Should not happen; forward-compat
				return nil, fmt.Errorf("unknown group %q", groupName)
			}
		}
	}
	return out, nil
}
