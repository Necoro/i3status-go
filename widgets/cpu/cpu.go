package cpu

import (
	u "github.com/Necoro/go-units"
	"github.com/prometheus/procfs"

	"github.com/Necoro/i3status-go/widgets"
)

const (
	name      = "cpu"
	icon rune = '\uf085' // 
)

type Params struct {
	Format string
}

type Widget struct {
	params    Params
	formatter widgets.Formatter
}

type Data struct {
	Frequencies []u.Value
}

func (w *Widget) Name() string {
	return name
}

func (w *Widget) Run() (d widgets.Data, err error) {
	d.Icon = icon

	proc, err := procfs.NewDefaultFS()
	if err != nil {
		return
	}

	cpuinfo, err := proc.CPUInfo()
	if err != nil {
		return
	}

	data := Data{
		Frequencies: make([]u.Value, len(cpuinfo)),
	}
	for i, cpu := range cpuinfo {
		data.Frequencies[i] = u.NewValue(cpu.CPUMHz, u.MegaHertz)
	}

	d.Text, err = w.format(data)
	return
}

func (w *Widget) format(data Data) (string, error) {
	if w.formatter == nil {
		if formatter, err := widgets.NewFormatter(w.params.Format); err != nil {
			return "", err
		} else {
			w.formatter = formatter
		}
	}

	return w.formatter(data)
}

func (w *Widget) Shutdown() {
}

func (w *Widget) Params() any {
	return &w.params
}

func init() {
	widgets.Register(name, func() widgets.Widget {
		return &Widget{params: Params{
			Format: "{{index .Frequencies 0}} {{index .Frequencies 1}}",
		}}
	})
}
