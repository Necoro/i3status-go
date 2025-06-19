package wifi

import (
	"net"

	u "github.com/Necoro/go-units"
	"github.com/mdlayher/wifi"

	"github.com/Necoro/i3status-go/widgets"
)

const name = "wifi"

type Params struct {
	// Format the output string. Takes a Go template string.
	Format    string
	Interface string
}

type Widget struct {
	params    Params
	client    *wifi.Client
	formatter widgets.Formatter
}

type Data struct {
	IPv6 string
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

func (w *Widget) Run() (d widgets.Data, err error) {
	if w.client == nil {
		if w.client, err = wifi.New(); err != nil {
			return
		}
	}
	ifs, err := w.client.Interfaces()
	if err != nil {
		return
	}

	var data Data
	for _, ifc := range ifs {
		if (w.params.Interface == "" && ifc.Name == "lo") ||
			ifc.Name == "" ||
			(w.params.Interface != "" && ifc.Name != w.params.Interface) {
			continue
		}

		data.Interface = ifc.Name

		var iface *net.Interface
		if iface, err = net.InterfaceByName(ifc.Name); err != nil {
			return
		}

		if net.FlagRunning&iface.Flags != 0 { // interface is up
			var lErr error
			bss, lErr := w.client.BSS(ifc)
			if lErr != nil {
				err = lErr
				return
			}

			data.SSID = bss.SSID
			data.Frequency = u.NewValue(float64(bss.Frequency), u.MegaHertz)

			addrs, lErr := iface.Addrs()
			if lErr != nil {
				err = lErr
				return
			}

			for _, addr := range addrs {
				ip, ok := addr.(*net.IPNet)
				if !ok || !ip.IP.IsGlobalUnicast() {
					continue
				}

				if ip.IP.To4() != nil {
					data.IPv4 = ip.String()
				} else {
					data.IPv6 = ip.String()
				}
			}
		}
		break
	}

	d.Text, err = w.format(data)
	return
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
