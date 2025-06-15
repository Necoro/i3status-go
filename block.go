package main

import (
	"fmt"

	"github.com/go-ini/ini"
	"github.com/go-viper/mapstructure/v2"

	"github.com/Necoro/i3status-go/widgets"
)

type Block struct {
	widgets.Widget `mapstructure:"-"`
	Interval       int
	ColorFg        string
}

var DefaultBlock = &Block{
	Interval: 5,
	ColorFg:  "#ffffff",
}

func NewBlock(section *ini.Section, defaults *Block) (*Block, error) {
	name := section.Name()
	var b Block
	if defaults != nil {
		b = *defaults
	}

	if name != ini.DefaultSection {
		w, err := widgets.Get(name)
		if err != nil {
			return nil, err
		}

		b.Widget = w
	}

	if err := b.loadValues(section.KeysHash()); err != nil {
		return nil, fmt.Errorf("failed to load values for block %s: %w", name, err)
	}

	return &b, nil
}

func (b *Block) loadValues(data map[string]string) error {
	metadata := mapstructure.Metadata{}
	config := &mapstructure.DecoderConfig{
		Metadata:         &metadata,
		WeaklyTypedInput: true,
		Result:           b,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err = decoder.Decode(data); err != nil {
		return err
	}

	if b.Widget == nil {
		// no widget given, no need to configure it
		return nil
	}

	// subset of keys not used for this block
	widgetData := make(map[string]string)
	for _, k := range metadata.Unused {
		widgetData[k] = data[k]
	}

	return widgets.Configure(b.Widget, widgetData)
}
