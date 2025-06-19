package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParams(t *testing.T) {
	config := `
foo = bar
fux = "baz"
a = 1.2
c = #ffffff
empty =
empty2 = ""
`
	expected := Config{
		GlobalParams: []Parameter{
			{"foo", "bar"},
			{"fux", "baz"},
			{"a", "1.2"},
			{"c", "#ffffff"},
			{"empty", ""},
			{"empty2", ""},
		},
	}

	cfg, err := cfgParser.ParseString("", config)

	assert.NoError(t, err)
	assert.Equal(t, expected, *cfg)
}

func TestSection(t *testing.T) {
	config := `
[section] 
foo = bar

[section.other]
fux = "baz"
a = 1.2

[ other ]
`
	expected := Config{
		Sections: []Section{
			{"section", "", []Parameter{
				{"foo", "bar"},
			}},
			{"section", "other", []Parameter{
				{"fux", "baz"},
				{"a", "1.2"},
			}},
			{"other", "", nil},
		},
	}

	cfg, err := cfgParser.ParseString("", config)

	assert.NoError(t, err)
	assert.Equal(t, expected, *cfg)
}

func TestComments(t *testing.T) {
	config := `
# a comment here
foo = bar

[section.other]
; another comment
fux = baz # not a comment

[comment] # comment also here
`
	expected := Config{
		Sections: []Section{
			{"section", "other", []Parameter{
				{"fux", "baz # not a comment"},
			}},
			{"comment", "", nil},
		},
		GlobalParams: []Parameter{
			{"foo", "bar"},
		},
	}

	cfg, err := cfgParser.ParseString("", config)

	assert.NoError(t, err)
	assert.Equal(t, expected, *cfg)
}

func TestInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			"trailing text with section",
			"[section] foo",
			"1:11: unexpected token \"foo\" (expected <eol> Parameter*)",
		},
		{
			"assignment directly with section",
			"[section] foo = bar",
			"1:11: unexpected token \"foo\" (expected <eol> Parameter*)",
		},
		{
			"multiple sections",
			"[section] [section2]",
			"1:11: unexpected token \"[\" (expected <eol> Parameter*)",
		},
		{
			"multiple strings",
			"foo = \"bar\" \"baz\"",
			"1:13: unexpected token \"baz\"",
		},
		{
			"inline comment",
			"foo #foo = bar",
			"1:15: unexpected token \"<EOF>\" (expected <assign> (<string> | <value>)?)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := cfgParser.ParseString("", test.input)
			assert.EqualError(t, err, test.err)
		})
	}
}
