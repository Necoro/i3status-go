package clock

import (
	"time"

	"golang.org/x/text/language"

	"github.com/goodsign/monday"

	"github.com/Necoro/i3status-go/widgets"
)

const name = "clock"

type Params struct {
	// Format is the Go time format string, see [time.Layout].
	Format string
	// Locale to use for things like names of days. If non-empty, we use [monday] to do the localization.
	// See [monday.Locale] for details and supported locales.
	Locale string
}

type Widget struct {
	params   Params
	formatFn func() string
}

func init() {
	widgets.Register(name, func() widgets.Widget {
		return &Widget{params: Params{
			Format: "Jan _2 Mon 15:04:05",
		}}
	})
}

func (c *Widget) Name() string {
	return name
}

func matchLocale(l string) monday.Locale {
	locales := monday.ListLocales()
	localeTags := make([]language.Tag, len(locales))
	for i := range locales {
		localeTags[i] = language.Make(string(locales[i]))
	}

	m := language.NewMatcher(localeTags)

	_, i, _ := m.Match(language.Make(l))
	return locales[i]
}

func (c *Widget) determineFormat() {
	formatStr := c.params.Format
	if c.params.Locale == "" {
		c.formatFn = func() string {
			return time.Now().Format(formatStr)
		}
	} else {
		locale := matchLocale(c.params.Locale)
		c.formatFn = func() string {
			return monday.Format(time.Now(), formatStr, locale)
		}
	}
}

func (c *Widget) Run() (widgets.Data, error) {
	if c.formatFn == nil {
		c.determineFormat()
	}

	return widgets.Data{
		Text: c.formatFn(),
	}, nil
}

func (c *Widget) Params() any {
	return &c.params
}

func (c *Widget) Shutdown() {
}
