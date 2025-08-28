package extract

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gi8lino/unveil/internal/quote"
	"github.com/gi8lino/unveil/internal/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractAll(t *testing.T) {
	t.Parallel()
	t.Run("Single JSON no quoting", func(t *testing.T) {
		dir := t.TempDir()
		jpath := filepath.Join(dir, "cfg.json")
		jcontent := `{"server":{"host":"localhost"}}`
		require.NoError(t, os.WriteFile(jpath, []byte(jcontent), 0o666))

		specs := []spec.ExtractSpec{{
			Kind:  spec.KindJSON,
			Path:  jpath,
			Key:   "server.host",
			Var:   "HOST",
			Quote: quote.QuoteNone,
		}}

		out, err := ExtractAll(specs)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"HOST": "localhost"}, out)
	})

	t.Run("Single YAML with single quotes", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		ypath := filepath.Join(dir, "cfg.yaml")
		ycontent := "foo: bar\n"
		require.NoError(t, os.WriteFile(ypath, []byte(ycontent), 0o666))

		specs := []spec.ExtractSpec{{
			Kind:  spec.KindYAML,
			Path:  ypath,
			Key:   "foo",
			Var:   "VAL",
			Quote: quote.QuoteSingle,
		}}

		out, err := ExtractAll(specs)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"VAL": "'bar'"}, out)
	})

	t.Run("Multiple specs mixed quoting", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		// file: key-value
		fpath := filepath.Join(dir, "app.env")
		require.NoError(t, os.WriteFile(fpath, []byte("K1=v1\n"), 0o666))

		// ini
		ipath := filepath.Join(dir, "conf.ini")
		icontent := "[S]\nK2=v2\n"
		require.NoError(t, os.WriteFile(ipath, []byte(icontent), 0o666))

		specs := []spec.ExtractSpec{
			{Kind: spec.KindFILE, Path: fpath, Key: "K1", Var: "V1", Quote: quote.QuoteNone},
			{Kind: spec.KindINI, Path: ipath, Key: "S.K2", Var: "V2", Quote: quote.QuoteDouble},
		}

		out, err := ExtractAll(specs)
		require.NoError(t, err)
		assert.Equal(t, "v1", out["V1"])
		assert.Equal(t, `"v2"`, out["V2"])
	})

	t.Run("TOML missing key error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		tpath := filepath.Join(dir, "cfg.toml")
		tcontent := ` [server] host = "h" `
		require.NoError(t, os.WriteFile(tpath, []byte(tcontent), 0o666))

		specs := []spec.ExtractSpec{{
			Kind:  spec.KindTOML,
			Path:  tpath,
			Key:   "server.port", // missing
			Var:   "PORT",
			Quote: quote.QuoteNone,
		}}

		out, err := ExtractAll(specs)
		require.Error(t, err)
		assert.Nil(t, out)
	})

	t.Run("JSON with json quoting", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		jpath := filepath.Join(dir, "special.json")
		jcontent := `{"msg":"line1\nline2\t\"q\"\\slash\\"}`
		require.NoError(t, os.WriteFile(jpath, []byte(jcontent), 0o666))

		specs := []spec.ExtractSpec{{
			Kind:  spec.KindJSON,
			Path:  jpath,
			Key:   "msg",
			Var:   "MSG",
			Quote: quote.QuoteJSON,
		}}

		out, err := ExtractAll(specs)
		require.NoError(t, err)
		// Expect a JSON-encoded string value (double-quoted with escapes)
		assert.Equal(t, `"line1\nline2\t\"q\"\\slash\\"`, out["MSG"])
	})
}
