package flag

import (
	"testing"

	"github.com/containeroo/tinyflags"
	"github.com/gi8lino/unveil/internal/quote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to fetch a dynamic group by name
func getGroup(t *testing.T, fs *tinyflags.FlagSet, name string) *tinyflags.DynamicGroup {
	t.Helper()
	for _, g := range fs.DynamicGroups() {
		if g.Name() == name {
			return g
		}
	}
	t.Fatalf("dynamic group %q not found", name)
	return nil
}

func TestParseFlags_Defaults_NoGroups(t *testing.T) {
	t.Parallel()

	t.Run("defaults", func(t *testing.T) {
		flags, err := ParseFlags([]string{}, "v", "c")
		require.NoError(t, err)
		require.NotNil(t, flags.FlagSet)
		assert.Equal(t, quote.QuoteNone, flags.Quote)
		// No groups instantiated
		assert.Len(t, flags.FlagSet.DynamicGroups(), 5) // json, yaml, file, toml, ini are registered
	})

	t.Run("global quote double", func(t *testing.T) {
		t.Parallel()

		flags, err := ParseFlags([]string{"--quote", "double"}, "v", "c")
		require.NoError(t, err)
		assert.Equal(t, quote.QuoteDouble, flags.Quote)
	})

	t.Run("JSON group single instance with override quote", func(t *testing.T) {
		t.Parallel()

		args := []string{
			"--json.web.path=./cfg.json",
			"--json.web.select=server.host",
			"--json.web.as=HOST",
			"--json.web.quote=json", // per-instance override
		}
		flags, err := ParseFlags(args, "v", "c")
		require.NoError(t, err)

		jsonGroup := getGroup(t, flags.FlagSet, "json")
		instances := jsonGroup.Instances()
		require.Len(t, instances, 1)
		id := instances[0]
		path := tinyflags.GetOrDefaultDynamic[string](jsonGroup, id, "path")
		sel := tinyflags.GetOrDefaultDynamic[string](jsonGroup, id, "select")
		as := tinyflags.GetOrDefaultDynamic[string](jsonGroup, id, "as")
		q := tinyflags.GetOrDefaultDynamic[string](jsonGroup, id, "quote")

		assert.Equal(t, "./cfg.json", path)
		assert.Equal(t, "server.host", sel)
		assert.Equal(t, "HOST", as)
		assert.Equal(t, "json", q)
	})

	t.Run("JSON group multiple instances with override quote", func(t *testing.T) {
		t.Parallel()

		// Missing --yaml.id.select (required) should cause Parse to error.
		args := []string{
			"--yaml.id.path=./cfg.yaml",
			// "--yaml.id.select" intentionally omitted
			"--yaml.id.as=FOO",
		}
		_, err := ParseFlags(args, "v", "c")
		require.Error(t, err)
	})

	t.Run("Multiple groups and instances", func(t *testing.T) {
		t.Parallel()

		args := []string{
			// json group: two instances
			"--json.a.path=./a.json",
			"--json.a.select=foo",
			"--json.a.as=FOO",
			"--json.b.path=./b.json",
			"--json.b.select=bar",
			"--json.b.as=BAR",

			// file group: one instance
			"--file.env.path=./app.env",
			"--file.env.select=API_KEY",
			"--file.env.as=API_KEY",
		}
		flags, err := ParseFlags(args, "v", "c")
		require.NoError(t, err)

		jsonGroup := getGroup(t, flags.FlagSet, "json")
		require.ElementsMatch(t, []string{"a", "b"}, jsonGroup.Instances())

		fileGroup := getGroup(t, flags.FlagSet, "file")
		require.ElementsMatch(t, []string{"env"}, fileGroup.Instances())

		// spot-check one instance field
		as := tinyflags.GetOrDefaultDynamic[string](fileGroup, "env", "as")
		assert.Equal(t, "API_KEY", as)
	})
}
