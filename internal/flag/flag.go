package flag

import (
	"github.com/containeroo/tinyflags"
	"github.com/gi8lino/unveil/internal/quote"
)

// Flags holds global options and the parsed FlagSet.
type Flags struct {
	Quote   quote.QuoteKind // global default quote mode
	Output  string          // output file
	Export  bool            // whether to export all variables
	FlagSet *tinyflags.FlagSet
}

// ParseFlags parses command-line arguments into Flags.
func ParseFlags(args []string, version, commit string) (Flags, error) {
	flags := Flags{}
	fs := tinyflags.NewFlagSet("unveil", tinyflags.ContinueOnError)

	fs.Version(version)
	fs.HelpText("show help")
	fs.VersionText("show version")

	// Global flags
	globalQuote := fs.String("quote", string(quote.QuoteNone), "quote all values").
		Choices(string(quote.QuoteNone), string(quote.QuoteSingle), string(quote.QuoteDouble), string(quote.QuoteJSON)).
		Placeholder("MODE").
		Value()
	fs.StringVar(&flags.Output, "output", "", "write to file instead of stdout").
		Value()
	fs.BoolVar(&flags.Export, "export", false, "add \"export\" prefix to all variables").
		Value()

	// Shared schema for dynamic groups
	registerGroup := func(name string) {
		g := fs.DynamicGroup(name)
		g.Title(name + " files:")
		g.String("path", "", "path to "+name+" file").Required()
		g.String("select", "", "selector/key to extract").
			Placeholder("KEY").
			Required()
		g.String("as", "", "destination variable. Defaults to ID in upper case.").
			Placeholder("VAR")
		g.String("quote", "", "quote mode").
			Choices(string(quote.QuoteNone), string(quote.QuoteSingle), string(quote.QuoteDouble), string(quote.QuoteJSON)).
			Placeholder("MODE")
	}

	for _, name := range []string{"json", "yaml", "file", "toml", "ini"} {
		registerGroup(name)
	}

	// Parse args
	if err := fs.Parse(args); err != nil {
		return Flags{}, err
	}
	flags.Quote = quote.QuoteKind(*globalQuote)
	flags.FlagSet = fs

	return flags, nil
}
