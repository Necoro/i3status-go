package main

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"

	"github.com/Necoro/i3status-go/widgets"
)

type Align string

const (
	AlignLeft   Align = "left"
	AlignCenter Align = "center"
	AlignRight  Align = "right"
)

type Markup string

const (
	MarkupPango Markup = "pango"
	MarkupNone  Markup = "none"
)

type Block struct {
	Widget    widgets.Widget `mapstructure:"-"`
	Qualifier string         `mapstructure:"-"`
	Interval  int
	ColorFg   string `mapstructure:"color"`
	Align     Align
	MinWidth  string
	Separator bool
	Markup    Markup
}

var defaultBlock = &Block{
	Interval:  5,
	ColorFg:   "#ffffff",
	Separator: true,
	Markup:    MarkupPango,
}

func GlobalBlock(globalParams Params) (*Block, error) {
	b := *defaultBlock
	err := b.loadValues(globalParams)
	return &b, err
}

func NewBlock(section Section, defaults *Block) (*Block, error) {
	var b Block
	if defaults != nil {
		b = *defaults
	}

	w, err := widgets.Get(section.Name)
	if err != nil {
		return nil, err
	}

	b.Widget = w
	b.Qualifier = section.Qualifier

	if err := b.loadValues(section.Params); err != nil {
		return nil, fmt.Errorf("failed to load values for block %s: %w", section.FullName(), err)
	}

	return &b, nil
}

func (b *Block) loadValues(data Params) error {
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
	widgetData := make(Params)
	for _, k := range metadata.Unused {
		widgetData[k] = data[k]
	}

	return widgets.Configure(b.Widget, widgetData)
}

func (b *Block) Run() widgets.Data {
	d := b.Widget.Run()

	if d.ColorFg == "" {
		d.ColorFg = b.ColorFg
	}

	return d
}
