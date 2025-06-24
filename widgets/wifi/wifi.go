package wifi

import (
	"bytes"
	"math"
	"net"

	u "github.com/Necoro/go-units"
	"github.com/mdlayher/wifi"

	"github.com/Necoro/i3status-go/widgets"
)

const (
	name      = "wifi"
	icon rune = '\uf1eb' // 
)

type Params struct {
	// Format the output string. Takes a Go template string.
	Format string
	// Interface to look for. If empty, it takes the first one.
	Interface string
	// DownFormat is the format for when the interface is down.
	DownFormat string
	// LevelGood is the lower bound for the good level.
	LevelGood int
	// LevelBad is the upper bound for the bad level.
	LevelBad int
}

type Widget struct {
	params        Params
	client        *wifi.Client
	formatter     widgets.Formatter[Data]
	downFormatter widgets.Formatter[Data]
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
	// Quality in percent (0-100).
	Quality u.Value
	up      bool
}

func (w *Widget) format(data Data) (string, error) {
	return w.formatter.Format(w.params.Format, data)
}

func (w *Widget) formatDown(data Data) (string, error) {
	return w.downFormatter.Format(w.params.DownFormat, data)
}

func (w *Widget) Name() string {
	return name
}

func computeSignalQuality(signal int) float64 {
	// from i3status-rust, in turn based on
	// <https://github.com/torvalds/linux/blob/9ff9b0d392ea08090cd1780fb196f36dbb586529/drivers/net/wireless/intel/ipw2x00/ipw2200.c#L4322-L4334>
	const (
		max  = -20.
		min  = -85.
		diff = max - min
	)

	sig64 := float64(signal)

	val := 100. - (max-sig64)*(15.*diff+62.*(max-sig64))/(diff*diff)
	return math.Max(0., math.Min(100., val))
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

			stationInfo, lErr := w.client.StationInfo(ifc)
			if lErr != nil {
				err = lErr
				return
			}
			for _, si := range stationInfo {
				if !bytes.Equal(si.HardwareAddr, bss.BSSID) {
					continue
				}
				qual := computeSignalQuality(si.SignalAverage)
				data.Quality = u.NewValue(qual, u.Percent)

				if int(qual) >= w.params.LevelGood {
					d.State = widgets.StateGood
				} else if int(qual) < w.params.LevelBad {
					d.State = widgets.StateBad
				} else {
					d.State = widgets.StateMid
				}
				break
			}

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
			LevelGood:  70,
			LevelBad:   20,
		}}
	})
}
