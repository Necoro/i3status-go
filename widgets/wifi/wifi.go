package wifi

import (
	"net"

	u "github.com/Necoro/go-units"
	"github.com/mdlayher/wifi"

	"github.com/Necoro/i3status-go/widgets"
)

const (
	name      = "wifi"
	icon rune = '\uf1eb'
)

type Params struct {
	// Format the output string. Takes a Go template string.
	Format string
	// Interface to look for. If empty, it takes the first one.
	Interface string
	// DownFormat is the format for when the interface is down.
	DownFormat string
}

type Widget struct {
	params        Params
	client        *wifi.Client
	formatter     widgets.Formatter
	downFormatter widgets.Formatter
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
	up        bool
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

func (w *Widget) formatDown(data Data) (string, error) {
	if w.downFormatter == nil {
		if downFormatter, err := widgets.NewFormatter(w.params.DownFormat); err != nil {
			return "", err
		} else {
			w.downFormatter = downFormatter
		}
	}

	return w.downFormatter(data)
}

func (w *Widget) Name() string {
	return name
}

func (w *Widget) Run() (d widgets.Data, err error) {
	d.Icon = icon

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

		data.up = net.FlagRunning&iface.Flags != 0
		if data.up {
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
					data.IPv4 = ip.IP.String()
				} else {
					data.IPv6 = ip.IP.String()
				}
			}
		}
		break
	}

	if data.up {
		d.Text, err = w.format(data)
	} else {
		d.Text, err = w.formatDown(data)
	}

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
			Format:     "{{.Interface}}: {{.SSID}} {{.Frequency | As GHz}}",
			DownFormat: "{{.Interface}}: not connected",
		}}
	})
}
