package widgets

import (
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

type Data struct {
	Text    string
	Urgent  bool
	Icon    rune
	ColorFg string
	ColorBg string
}

func (d Data) FullText() string {
	if d.Icon != 0 {
		return string(d.Icon) + " " + d.Text
	}
	return d.Text
}

type Widget interface {
	// Name of this widget
	Name() string
	// Run the widget
	Run() (Data, error)
	// Shutdown the widget cleanly
	Shutdown()
	// Params of the widget to be filled with configured values.
	// Must return a map[string] or a pointer to a struct.
	// The returned value can already contain default values.
	Params() any
}

type WidgetConstructor = func() Widget

var registry = make(map[string]WidgetConstructor)

func Register(name string, w WidgetConstructor) {
	name = strings.ToLower(name)
	if _, ok := registry[name]; ok {
		panic("Duplicate widget: " + name)
	}

	registry[name] = w
}

func Get(name string) (Widget, error) {
	name = strings.ToLower(name)
	w, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown widget: %s", name)
	}
	return w(), nil
}

func Configure(w Widget, data map[string]string, config *mapstructure.DecoderConfig) error {
	p := w.Params()
	config.Result = p

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create decoder for %s: %w", w.Name(), err)
	}

	if err = decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode values for %s: %w", w.Name(), err)
	}

	return nil
}
