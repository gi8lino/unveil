package quote

import (
	"encoding/json"
	"strings"
)

// QuoteKind defines how to render the VALUE (not the key).
type QuoteKind string

const (
	QuoteNone   QuoteKind = "none"   // write raw value
	QuoteSingle QuoteKind = "single" // '…' POSIX style
	QuoteDouble QuoteKind = "double" // "…" with escaping
	QuoteJSON   QuoteKind = "json"   // JSON-encode string
)

// QuoteValue applies the quoting/encoding to v.
func QuoteValue(v string, kind QuoteKind) string {
	switch kind {
	case QuoteNone:
		return v
	case QuoteSingle:
		if v == "" {
			return "''"
		}
		// Escape embedded single quotes: abc'def → 'abc'\''def'
		return "'" + strings.ReplaceAll(v, "'", `'\''`) + "'"
	case QuoteDouble:
		if v == "" {
			return `""`
		}
		r := v
		r = strings.ReplaceAll(r, `\`, `\\`)
		r = strings.ReplaceAll(r, `"`, `\"`)
		r = strings.ReplaceAll(r, `$`, `\$`)
		r = strings.ReplaceAll(r, "`", "\\`")
		return `"` + r + `"`
	case QuoteJSON:
		b, _ := json.Marshal(v)
		return string(b)
	default:
		return v
	}
}
