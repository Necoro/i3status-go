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

type Formatter func(data any) (string, error)

func NewFormatter(format string) (Formatter, error) {
	tpl, err := template.New("Format").
		Funcs(template.FuncMap{
			"Scale": scale,
			"As":    as,
		}).
		Parse(format)

	if err != nil {
		return nil, fmt.Errorf("parsing Format: %w", err)
	}

	return func(data any) (string, error) {
		var sb strings.Builder
		if err := tpl.Execute(&sb, data); err != nil {
			return "", fmt.Errorf("formatting: %w", err)
		}

		return sb.String(), nil
	}, nil
}

// scale version needed for the template
func scale(digits int, value u.Value) string {
	hValue := value.Humanize()
	return hValue.Fmt(u.FmtOptions{
		Short:     true,
		Precision: digits,
		Label:     true,
	})
}

// as converts the unit
func as(symbol string, value u.Value) u.Value {
	unit, _ := u.Find(symbol)
	if unit == nil {
		return value
	}
	return value.MustConvert(unit)
}
