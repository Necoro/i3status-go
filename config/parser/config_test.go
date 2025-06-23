package parser

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestSimpleParams(t *testing.T) {
	config := `
foo = bar
Fux = "baz"
a = 1.2
c = #FFffff
empty =
empty2 = ""
`
	expected := &Config{
		Global: Params{
			"foo":    "bar",
			"fux":    "baz",
			"a":      "1.2",
			"c":      "#FFffff",
			"empty":  "",
			"empty2": "",
		},
		Sections: []Section{},
	}

	cfg, err := Parse("", []byte(config))

	assert.NoError(t, err)
	assert.Equal(t, expected, cfg.(*Config))
}

func TestDuplicateParams(t *testing.T) {
	config := `
foo = bar
fux = "baz"
capC = minC
foo = 1.2
fux =
capc = none
`
	expected := &Config{
		Global: Params{
			"foo":  "1.2",
			"fux":  "",
			"capc": "none",
		},
		Sections: []Section{},
	}

	cfg, err := Parse("", []byte(config))

	assert.NoError(t, err)
	assert.Equal(t, expected, cfg.(*Config))
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
	expected := &Config{
		Global: Params{},
		Sections: []Section{
			{"section", "", Params{
				"foo": "bar",
			}},
			{"section", "other", Params{
				"fux": "baz",
				"a":   "1.2",
			}},
			{"other", "", Params{}},
		},
	}

	cfg, err := Parse("", []byte(config))

	assert.NoError(t, err)
	assert.Equal(t, expected, cfg.(*Config))
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
	expected := &Config{
		Sections: []Section{
			{"section", "other", Params{
				"fux": "baz # not a comment",
			}},
			{"comment", "", Params{}},
		},
		Global: Params{
			"foo": "bar",
		},
	}

	cfg, err := Parse("", []byte(config))

	assert.NoError(t, err)
	assert.Equal(t, expected, cfg.(*Config))
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
			"1:11 (10): no match found, expected: \"\\n\", [ \\t] or [;#]",
		},
		{
			"assignment directly with section",
			"[section] foo = bar",
			"1:11 (10): no match found, expected: \"\\n\", [ \\t] or [;#]",
		},
		{
			"multiple sections",
			"[section] [section2]",
			"1:11 (10): no match found, expected: \"\\n\", [ \\t] or [;#]",
		},
		{
			"multiple strings",
			"foo = \"bar\" \"baz\"",
			"1:13 (12): no match found, expected: \"\\n\" or [ \\t]",
		},
		{
			"inline comment",
			"foo #foo = bar",
			"1:5 (4): no match found, expected: [ \\t] or [=:]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse("", []byte(test.input))
			assert.EqualError(t, err, test.err)
		})
	}
}
