package widgets

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
)

type Widget interface {
	// Name of this widget
	Name() string
	// Run the widget
	Run()
	// Params of the widget to be filled with configured values.
	// Must return a map[string] or a pointer to a struct.
	// The returned value can already contain default values.
	Params() any
}

type WidgetConstructor = func() Widget

var registry = make(map[string]WidgetConstructor)

func Register(name string, w WidgetConstructor) {
	if _, ok := registry[name]; ok {
		panic("Duplicate widget: " + name)
	}

	registry[name] = w
}

func Configure(w Widget, data map[string]string) error {
	p := w.Params()

	metadata := mapstructure.Metadata{}
	config := &mapstructure.DecoderConfig{
		Metadata:         &metadata,
		WeaklyTypedInput: true,
		Result:           p,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create decoder for %s: %w", w.Name(), err)
	}

	if err = decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode values for %s: %w", w.Name(), err)
	}

	return nil
}
