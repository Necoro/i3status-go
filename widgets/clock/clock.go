package clock

import "github.com/Necoro/i3status-go/widgets"

const name = "clock"

type Params struct {
	Format string
	Locale string
}

type Widget struct {
	params Params
}

func init() {
	widgets.Register(name, func() widgets.Widget {
		return &Widget{}
	})
}

func (c *Widget) Name() string {
	return name
}

func (c *Widget) Run() {
	panic("implement me")
}

func (c *Widget) Params() any {
	return &Params{Format: "Jan _2 Mon 15:04:05"}
}
