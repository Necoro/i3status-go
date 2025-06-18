package wifi

import (
	u "github.com/Necoro/go-units"
	"github.com/mdlayher/wifi"

	"github.com/Necoro/i3status-go/widgets"
)

const name = "wifi"

type Params struct {
	//Format the output string. Takes a Go template string.
	Format    string
	Interface string
}

type Widget struct {
	params    Params
	client    *wifi.Client
	formatter widgets.Formatter
}

type Data struct {
	IPv4 string
	// SSID of the connected network.
	SSID string
	// Interface name, e.g. wlan0, wlan1, etc.
	Interface string
	// Frequency in MHz
	Frequency u.Value
	Signal    string
	Quality   string
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

func (w *Widget) Name() string {
	return name
}

func (w *Widget) Run() widgets.Data {
	if w.client == nil {
		w.client, _ = wifi.New()
	}
	ifs, _ := w.client.Interfaces()

	var data Data
	for _, ifc := range ifs {
		if (w.params.Interface == "" && ifc.Name == "lo") ||
			ifc.Name == "" ||
			(w.params.Interface != "" && ifc.Name != w.params.Interface) {
			continue
		}
		data.Interface = ifc.Name

		bss, _ := w.client.BSS(ifc)
		data.SSID = bss.SSID
		data.Frequency = u.NewValue(float64(bss.Frequency), u.MegaHertz)
		break
	}

	text, err := w.format(data)
	if err != nil {
		panic(err)
	}
	return widgets.Data{
		Text: text,
	}

}

func (w *Widget) Shutdown() {
	if w.client != nil {
		_ = w.client.Close()
		w.client = nil
	}
}

func (w *Widget) Params() any {
	return &w.params
}

func init() {
	widgets.Register(name, func() widgets.Widget {
		return &Widget{params: Params{
			Format: "{{.Interface}}: {{.SSID}} {{.Frequency | As GHz}}",
		}}
	})
}
