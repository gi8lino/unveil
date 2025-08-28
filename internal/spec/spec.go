package spec

import "github.com/gi8lino/unveil/internal/quote"

// Kind is the file/source type.
type Kind string

const (
	KindJSON Kind = "json"
	KindYAML Kind = "yaml"
	KindTOML Kind = "toml"
	KindINI  Kind = "ini"
	KindFILE Kind = "file"
)

// ExtractSpec describes one extraction instruction.
type ExtractSpec struct {
	Kind  Kind            // source type
	Path  string          // file path
	Key   string          // selector/key/path inside file
	Var   string          // destination env var name
	Quote quote.QuoteKind // quoting mode for value
}
