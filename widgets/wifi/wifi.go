package wifi

import (
	"fmt"

	"github.com/mdlayher/wifi"

	"github.com/Necoro/i3status-go/widgets"
)

const name = "wifi"

type Params struct {
	Format    string
	Interface string
}

type Widget struct {
	params   Params
	formatFn func() string
	client   *wifi.Client
}

type Data struct {
	Interface string
	Frequency int
	Signal    string
	Quality   string
	SSID      string
}

func (w Widget) Name() string {
	return name
}

func (w Widget) Run() widgets.Data {
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
		data.Frequency = bss.Frequency
		break
	}

	return widgets.Data{
		Text: fmt.Sprintf("%s: %s %d", data.Interface, data.SSID, data.Frequency),
	}

}

func (w Widget) Shutdown() {
	if w.client != nil {
		_ = w.client.Close()
		w.client = nil
	}
}

func (w Widget) Params() any {
	return &w.params
}

func init() {
	widgets.Register(name, func() widgets.Widget {
		return &Widget{params: Params{
			Format: "%ip %essid",
		}}
	})
}
