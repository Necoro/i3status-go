package main

import (
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"

	"github.com/Necoro/i3status-go/config"
	"github.com/Necoro/i3status-go/widgets"
)

const errorColor = "#FF262A"

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

type BlockConfig struct {
	Interval  int
	ColorFg   string `config:"color"`
	ColorGood string
	ColorMid  string
	ColorBad  string
	Align     Align
	MinWidth  string
	Separator bool
	Markup    Markup
	Icon      rune
}

type Block struct {
	Widget      widgets.Widget `config:"-"`
	Qualifier   string         `config:"-"`
	BlockConfig `config:",squash"`
}

var defaultConfig = BlockConfig{
	Interval:  5,
	ColorFg:   "#ffffff",
	ColorGood: "#ffffff",
	ColorMid:  "#ffff00",
	ColorBad:  errorColor,
	Separator: true,
	Markup:    MarkupNone,
}

func GlobalBlock(globalParams config.Params) (*Block, error) {
	b := &Block{BlockConfig: defaultConfig}
	err := b.loadValues(globalParams)
	return b, err
}

func NewBlock(section config.Section, defaults *Block) (*Block, error) {
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

func matchConfigKey(mapKey, fieldName string) bool {
	mapKey = strings.ReplaceAll(mapKey, "_", "")
	return strings.EqualFold(mapKey, fieldName)
}

func decoderConfig() *mapstructure.DecoderConfig {
	metadata := mapstructure.Metadata{}
	return &mapstructure.DecoderConfig{
		Metadata:         &metadata,
		WeaklyTypedInput: true,
		TagName:          "config",
		MatchName:        matchConfigKey,
	}
}

func (b *Block) loadValues(data config.Params) error {
	decCfg := decoderConfig()
	decCfg.Result = b

	decoder, err := mapstructure.NewDecoder(decCfg)
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
	widgetData := make(config.Params)
	for _, k := range decCfg.Metadata.Unused {
		widgetData[k] = data[k]
	}

	return widgets.Configure(b.Widget, widgetData, decoderConfig())
}

func (b *Block) Run() widgets.Data {
	d, err := b.Widget.Run()

	if err != nil {
		return widgets.Data{
			Text:    err.Error(),
			Urgent:  true,
			ColorFg: errorColor,
		}
	}

	if d.ColorFg == "" {
		switch d.State {
		case widgets.StateGood:
			d.ColorFg = b.ColorGood
		case widgets.StateMid:
			d.ColorFg = b.ColorMid
		case widgets.StateBad:
			d.ColorFg = b.ColorBad
		default:
			d.ColorFg = b.ColorFg
		}
	}

	if b.Icon != 0 {
		d.Icon = b.Icon
	}

	return d
}
