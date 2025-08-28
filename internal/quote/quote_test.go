package quote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuoteValue_None_Raw(t *testing.T) {
	t.Parallel()

	t.Run("None Raw", func(t *testing.T) {
		t.Parallel()
		in := `plain value`
		out := QuoteValue(in, QuoteNone)
		assert.Equal(t, in, out)
	})

	t.Run("Single Empty", func(t *testing.T) {
		t.Parallel()
		out := QuoteValue("", QuoteSingle)
		assert.Equal(t, "''", out)
	})

	t.Run("Single simple", func(t *testing.T) {
		t.Parallel()
		out := QuoteValue("hello world", QuoteSingle)
		assert.Equal(t, "'hello world'", out)
	})

	t.Run("Single Escape Quote", func(t *testing.T) {
		t.Parallel()
		out := QuoteValue("abc'def", QuoteSingle)
		// POSIX-safe single-quote escaping: 'abc'\''def'
		assert.Equal(t, "'abc'\\''def'", out)
	})

	t.Run("Double Empty", func(t *testing.T) {
		t.Parallel()
		out := QuoteValue("", QuoteDouble)
		assert.Equal(t, `""`, out)
	})

	t.Run("Escape Backslash Quote dollar Backtick", func(t *testing.T) {
		t.Parallel()
		in := `a\b"c$d` + "`" + `e`
		out := QuoteValue(in, QuoteDouble)
		require.True(t, len(out) >= 2 && out[0] == '"' && out[len(out)-1] == '"', "must be wrapped in double quotes")

		// Inside should have escapes for \ " $ `
		expectedInner := `a\\b\"c\$d\` + "`" + `e`
		assert.Equal(t, `"`+expectedInner+`"`, out)
	})

	t.Run("Simple JSON", func(t *testing.T) {
		t.Parallel()
		out := QuoteValue("hello", QuoteJSON)
		assert.Equal(t, `"hello"`, out)
	})

	t.Run("JSON With Specials", func(t *testing.T) {
		t.Parallel()
		in := "line1\nline2\t\"q\"\\slash\\"
		out := QuoteValue(in, QuoteJSON)
		// JSON-encoded form must be quoted and contain escapes for \n, \t, \" and \\.
		assert.Equal(t, `"line1\nline2\t\"q\"\\slash\\"`, out)
	})

	t.Run("Unknown Kind fallback", func(t *testing.T) {
		t.Parallel()
		unknown := QuoteKind("weird")
		in := `keep me`
		out := QuoteValue(in, unknown)
		assert.Equal(t, in, out)
	})
}
