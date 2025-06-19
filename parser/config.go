// noinspection GoVetStructTagInspection
package parser

import (
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Config struct {
	GlobalParams []Parameter `@@*`
	Sections     []Section   `@@*`
}

type Parameter struct {
	Name  string `@Ident Assign`
	Value string `(@String | @Value)?`
}

type Section struct {
	Name      string      `"[" @Ident`
	Qualifier string      `("." @Ident)? "]" EOL`
	Params    []Parameter `@@*`
}

var (
	cfgLexer = lexer.MustStateful(lexer.Rules{
		"Root": {
			{"Ident", `[a-zA-Z_][a-zA-Z_0-9]*`, nil},
			{"comment", `[#;][^\n]+`, nil},
			{"whitespace", `[ \t]+`, nil},
			{"eol", `\n`, nil},
			{"Assign", `[:=]`, lexer.Push("Assign")},
			{"SecStart", `\[`, lexer.Push("Section")},
		},
		"Assign": {
			{"eol", `\n`, lexer.Pop()},
			{"whitespace", `[ \t]+`, nil},
			{"String", `"[^"]*"`, nil},
			{"Value", `[^\n]+`, nil},
		},
		"Section": {
			{"EOL", `\n`, lexer.Pop()},
			lexer.Include("Root"),
			{"Punctuation", `\[|\]|\.`, nil},
		},
	})
	cfgParser = participle.MustBuild[Config](
		participle.Lexer(cfgLexer),
		participle.Unquote("String"),
	)
)

func Parse(filename string, r io.Reader) (*Config, error) {
	return cfgParser.Parse(filename, r)
}
