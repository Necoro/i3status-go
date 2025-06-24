package widgets

import (
	"fmt"
	"strings"
	"text/template"

	u "github.com/Necoro/go-units"
)

func init() {
	u.DefaultFmtOptions.Short = true
}

type Formatter[D any] struct {
	fn func(data D) (string, error)
}

func (f *Formatter[D]) build(format string) error {
	tpl, err := template.New("Format").
		Funcs(template.FuncMap{
			"Scale":  scale,
			"As":     as,
			"NoUnit": noUnit,
		}).
		Parse(format)

	if err != nil {
		return fmt.Errorf("parsing Format: %w", err)
	}

	f.fn = func(data D) (string, error) {
		var sb strings.Builder
		if err := tpl.Execute(&sb, data); err != nil {
			return "", fmt.Errorf("formatting: %w", err)
		}

		return sb.String(), nil
	}
	return nil
}

func (f *Formatter[D]) Format(format string, data D) (string, error) {
	if f.fn == nil {
		if err := f.build(format); err != nil {
			return "", err
		}
	}

	return f.fn(data)
}

type valueFormat struct {
	value u.Value
	opts  u.FmtOptions
}

func newValueFormat(value u.Value) valueFormat {
	return valueFormat{value: value, opts: u.FmtOptions{
		Short:     true,
		Precision: -1,
		Label:     true,
	}}
}

func (f valueFormat) String() string {
	return f.value.Fmt(f.opts)
}

func valueFormatter(data any, fn func(vf valueFormat) valueFormat) valueFormat {
	switch data := data.(type) {
	case u.Value:
		return fn(newValueFormat(data))
	case valueFormat:
		return fn(data)
	default:
		panic("invalid type")
	}
}

// scale version needed for the template
func scale(digits int, data any) valueFormat {
	return valueFormatter(data, func(vf valueFormat) valueFormat {
		vf.opts.Precision = digits
		return vf
	})
}

func noUnit(data any) valueFormat {
	return valueFormatter(data, func(vf valueFormat) valueFormat {
		vf.opts.Label = false
		return vf
	})
}

// as converts the unit
func as(symbol string, value u.Value) u.Value {
	unit, _ := u.Find(symbol)
	if unit == nil || value.Unit() == nil {
		return value
	}
	return value.MustConvert(unit)
}
