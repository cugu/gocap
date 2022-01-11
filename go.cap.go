package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type File struct {
	Self     Line   `parser:"@@"`
	External []Line `parser:"@@*"`

	expected map[string][]string
}

func (f *File) setup() {
	expected := map[string][]string{
		f.Self.Package: f.Self.Capabilities,
	}
	for _, ext := range f.External {
		expected[ext.Package] = ext.Capabilities
	}
	f.expected = expected
}

type Line struct {
	Package      string   `parser:"@String"`
	Capabilities []string `parser:" '(' ((@String ',' )* @String)? ')'"`
}

var parser = participle.MustBuild(&File{},
	participle.Lexer(lexer.MustSimple([]lexer.Rule{
		{Name: "Comment", Pattern: `(?:#|//)[^\n]*\n?`},
		{Name: "String", Pattern: `[^\s,()]+`},
		{Name: "Punct", Pattern: `[,()]`},
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	})),
	participle.Elide("Comment", "Whitespace"),
	participle.UseLookahead(2),
)

func parseGoCap() (*File, error) {
	r, err := os.Open("go.cap")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New("go.cap file does not exist")
		}
		return nil, fmt.Errorf("could not open go.cap file: %s", err)
	}
	defer r.Close()

	return parse(r)
}

func parse(r io.Reader) (*File, error) {
	ast := &File{}
	err := parser.Parse("", r, ast)
	if err != nil {
		return nil, err
	}

	ast.setup()

	return ast, err
}
