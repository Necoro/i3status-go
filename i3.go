package main

import "github.com/Necoro/i3status-go/widgets"

// I3BarHeader represents the header of an i3bar message.
type I3BarHeader struct {
	Version     uint8 `json:"version"`
	StopSignal  *int  `json:"stop_signal,omitempty"`
	ContSignal  *int  `json:"cont_signal,omitempty"`
	ClickEvents bool  `json:"click_events,omitempty"`
}

// I3BarBlock represents a block of i3bar message.
type I3BarBlock struct {
	FullText            string  `json:"full_text,omitempty"`
	ShortText           string  `json:"short_text,omitempty"`
	Color               string  `json:"color,omitempty"`
	BorderColor         string  `json:"border,omitempty"`
	BorderTop           *uint16 `json:"border_top,omitempty"`
	BorderBottom        *uint16 `json:"border_bottom,omitempty"`
	BorderLeft          *uint16 `json:"border_left,omitempty"`
	BorderRight         *uint16 `json:"border_right,omitempty"`
	BackgroundColor     string  `json:"background,omitempty"`
	Markup              string  `json:"markup,omitempty"`
	MinWidth            *uint16 `json:"min_width,omitempty"`
	Align               string  `json:"align,omitempty"`
	Name                string  `json:"name,omitempty"`
	Instance            string  `json:"instance,omitempty"`
	Urgent              bool    `json:"urgent,omitempty"`
	Separator           *bool   `json:"separator,omitempty"`
	SeparatorBlockWidth uint16  `json:"separator_block_width,omitempty"`
}

func NewI3BarBlock(d widgets.Data) I3BarBlock {
	return I3BarBlock{
		FullText:        d.FullText(),
		Color:           d.ColorFg,
		BackgroundColor: d.ColorBg,
		Urgent:          d.Urgent,
	}
}

// I3BarClickEvent represents a user click event message.
type I3BarClickEvent struct {
	Name      string   `json:"name,omitempty"`
	Instance  string   `json:"instance,omitempty"`
	Button    uint8    `json:"button"`
	X         uint16   `json:"x"`
	Y         uint16   `json:"y"`
	RelativeX uint16   `json:"relative_x"`
	RelativeY uint16   `json:"relative_y"`
	OutputX   uint16   `json:"output_x"`
	OutputY   uint16   `json:"output_y"`
	Width     uint16   `json:"width"`
	Height    uint16   `json:"height"`
	Modifiers []string `json:"modifiers"`
}
