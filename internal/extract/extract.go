package extract

import (
	"fmt"

	"github.com/gi8lino/unveil/internal/quote"
	"github.com/gi8lino/unveil/internal/spec"

	"github.com/containeroo/resolver"
)

// ExtractAll resolves all specs and returns a map of KEY=VALUE.
// Values are quoted according to spec.Quote.
func ExtractAll(specs []spec.ExtractSpec) (map[string]string, error) {
	out := make(map[string]string, len(specs))
	for _, s := range specs {
		// Compose resolver URI: "kind:/path//key"
		value := fmt.Sprintf("%s:%s//%s", s.Kind, s.Path, s.Key)
		val, err := resolver.ResolveVariable(value)
		if err != nil {
			return nil, fmt.Errorf("%s %q (%s=%q): %w", s.Kind, s.Var, "path", s.Path, err)
		}
		// Apply quoting policy
		out[s.Var] = quote.QuoteValue(val, s.Quote)
	}
	return out, nil
}
