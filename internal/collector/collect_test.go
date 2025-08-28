package collector

import (
	"testing"

	"github.com/gi8lino/unveil/internal/flag"
	"github.com/gi8lino/unveil/internal/quote"
	"github.com/gi8lino/unveil/internal/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func byVar(specs []spec.ExtractSpec) map[string]spec.ExtractSpec {
	out := make(map[string]spec.ExtractSpec, len(specs))
	for _, s := range specs {
		out[s.Var] = s
	}
	return out
}

func TestCollect(t *testing.T) {
	t.Parallel()

	t.Run("JSON instance; default var from id; per-instance quote overrides global", func(t *testing.T) {
		t.Parallel()

		args := []string{
			"--quote", "single", // global default
			"--json.web.path=/cfg.json",
			"--json.web.select=server.host",
			"--json.web.quote=double", // override
		}
		flags, err := flag.ParseFlags(args, "v", "c")
		require.NoError(t, err)

		got, err := Collect(&flags)
		require.NoError(t, err)
		require.Len(t, got, 1)

		s := got[0]
		assert.Equal(t, spec.KindJSON, s.Kind)
		assert.Equal(t, "/cfg.json", s.Path)
		assert.Equal(t, "server.host", s.Key)
		assert.Equal(t, "WEB", s.Var)               // derived from id: "web" → "WEB"
		assert.Equal(t, quote.QuoteDouble, s.Quote) // per-instance override
	})

	t.Run("YAML instance; explicit var; inherits global quote when no override", func(t *testing.T) {
		t.Parallel()

		args := []string{
			"--quote", "json", // global
			"--yaml.app.path=/cfg.yaml",
			"--yaml.app.select=foo",
			"--yaml.app.as=FOO",
		}
		flags, err := flag.ParseFlags(args, "v", "c")
		require.NoError(t, err)

		got, err := Collect(&flags)
		require.NoError(t, err)
		require.Len(t, got, 1)

		s := got[0]
		assert.Equal(t, spec.KindYAML, s.Kind)
		assert.Equal(t, "/cfg.yaml", s.Path)
		assert.Equal(t, "foo", s.Key)
		assert.Equal(t, "FOO", s.Var)
		assert.Equal(t, quote.QuoteJSON, s.Quote) // inherited from global
	})

	t.Run("TOML instance; explicit var; inherits global quote when no override", func(t *testing.T) {
		t.Parallel()

		args := []string{
			"--quote", "json", // global
			"--toml.app.path=/cfg.toml",
			"--toml.app.select=foo",
			"--toml.app.as=FOO",
			"--toml.app.quote=double", // override
		}
		flags, err := flag.ParseFlags(args, "v", "c")
		require.NoError(t, err)

		got, err := Collect(&flags)
		require.NoError(t, err)
		require.Len(t, got, 1)

		s := got[0]
		assert.Equal(t, spec.KindTOML, s.Kind)
		assert.Equal(t, "/cfg.toml", s.Path)
		assert.Equal(t, "foo", s.Key)
		assert.Equal(t, "FOO", s.Var)
		assert.Equal(t, quote.QuoteDouble, s.Quote) // inherited from global (default)
	})

	t.Run("Multiple groups and instances", func(t *testing.T) {
		t.Parallel()

		args := []string{
			// json group: two instances
			"--json.a.path=/a.json",
			"--json.a.select=foo",
			"--json.a.as=FOO",
			"--json.b.path=/b.json",
			"--json.b.select=bar",
			// no 'as' → should become "B"

			// file group: one instance with override quote
			"--file.env.path=/app.env",
			"--file.env.select=API_KEY",
			"--file.env.as=KEY",
			"--file.env.quote=single",

			// ini group: one instance
			"--ini.db.path=/conf.ini",
			"--ini.db.select=Section.User",
			"--ini.db.as=DBUSER",

			// global quote should be 'none' by default
		}
		flags, err := flag.ParseFlags(args, "v", "c")
		require.NoError(t, err)

		specs, err := Collect(&flags)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(specs), 4)

		m := byVar(specs)

		// json.a → explicit var and default quote (none)
		foo, ok := m["FOO"]
		require.True(t, ok)
		assert.Equal(t, spec.KindJSON, foo.Kind)
		assert.Equal(t, "/a.json", foo.Path)
		assert.Equal(t, "foo", foo.Key)
		assert.Equal(t, quote.QuoteNone, foo.Quote)

		// json.b → var from id "b" → "B"
		b, ok := m["B"]
		require.True(t, ok)
		assert.Equal(t, spec.KindJSON, b.Kind)
		assert.Equal(t, "/b.json", b.Path)
		assert.Equal(t, "bar", b.Key)
		assert.Equal(t, quote.QuoteNone, b.Quote)

		// file.env → explicit var + per-instance quote override
		key, ok := m["KEY"]
		require.True(t, ok)
		assert.Equal(t, spec.KindFILE, key.Kind)
		assert.Equal(t, "/app.env", key.Path)
		assert.Equal(t, "API_KEY", key.Key)
		assert.Equal(t, quote.QuoteSingle, key.Quote)

		// ini.db → explicit var, default quote
		dbu, ok := m["DBUSER"]
		require.True(t, ok)
		assert.Equal(t, spec.KindINI, dbu.Kind)
		assert.Equal(t, "/conf.ini", dbu.Path)
		assert.Equal(t, "Section.User", dbu.Key)
		assert.Equal(t, quote.QuoteNone, dbu.Quote)
	})
}
