package mem

import (
	u "github.com/Necoro/go-units"
	"github.com/Necoro/i3status-go/widgets"
	"github.com/prometheus/procfs"
)

const (
	name      = "mem"
	icon rune = '\uf2db' // 
)

type Params struct {
	// Format the output string. Takes a Go template string.
	Format string
	// LevelGood is the upper bound of 'used percentage' for the good level.
	LevelGood int
	// LevelBad is the lower bound of 'used percentage' for the bad level.
	LevelBad int
}

type Widget struct {
	params    Params
	formatter widgets.Formatter[Data]
	proc      *procfs.FS
}

type Data struct {
	// Total memory
	Total u.Value
	// Used memory. Calculated as Total - Available
	// This is *not* identical to how htop calculates it.
	Used        u.Value
	UsedPercent u.Value
	// Free memory. This is not identical to Available.
	// See https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=34e431b0ae398fc54ea69ff85ec700722c9da773
	Free        u.Value
	FreePercent u.Value
	// Available memory. This is Free plus some buffers and caches that can be freed, if so needed.
	// This is directly reported by the kernel.
	Available        u.Value
	AvailablePercent u.Value
	Swap             struct {
		// Total swap.
		Total u.Value
		// Free swap.
		Free        u.Value
		FreePercent u.Value
		// Used swap. This may be less than Total - Free, because there exists also cached swap.
		// That is swap that has been loaded into memory but not yet removed from swap.
		Used        u.Value
		UsedPercent u.Value
	}
}

func val(v *uint64) u.Value {
	if v == nil {
		return u.NewValue(0, u.Kibibyte)
	}

	return u.NewValue(float64(*v), u.Kibibyte)
}

func perc(num, denom *uint64) u.Value {
	if num == nil || denom == nil || *num == 0 || *denom == 0 {
		return u.NewValue(0, u.Percent)
	}

	return u.NewValue(float64(*num)/float64(*denom)*100, u.Percent)
}

func (w *Widget) readData() (d Data, err error) {
	if w.proc == nil {
		fs, _err := procfs.NewDefaultFS()
		if _err != nil {
			err = _err
		}
		w.proc = &fs
	}

	mem, err := w.proc.Meminfo()
	if err != nil {
		return
	}

	d.Total = val(mem.MemTotal)
	d.Free = val(mem.MemFree)
	d.FreePercent = perc(mem.MemFree, mem.MemTotal)
	d.Available = val(mem.MemAvailable)
	d.AvailablePercent = perc(mem.MemAvailable, mem.MemTotal)

	d.Used = u.NewValue(d.Total.Float()-d.Available.Float(), u.Kibibyte)
	d.UsedPercent = u.NewValue(100.-d.AvailablePercent.Float(), u.Percent)

	d.Swap.Total = val(mem.SwapTotal)
	d.Swap.Free = val(mem.SwapFree)
	d.Swap.FreePercent = perc(mem.SwapFree, mem.SwapTotal)

	swapCached := val(mem.SwapCached)
	swapUsed := uint64(d.Swap.Total.Float() - d.Swap.Free.Float() - swapCached.Float())
	d.Swap.Used = val(&swapUsed)
	d.Swap.UsedPercent = perc(&swapUsed, mem.SwapTotal)

	return
}

func (w *Widget) Run() (d widgets.Data, err error) {
	d.Icon = icon

	data, err := w.readData()
	if err != nil {
		return
	}

	d.Text, err = w.format(data)

	usedPercent := int(data.UsedPercent.Float())
	if usedPercent < w.params.LevelGood {
		d.State = widgets.StateGood
	} else if usedPercent >= w.params.LevelBad {
		d.State = widgets.StateBad
	} else {
		d.State = widgets.StateMid
	}
	return
}

func (w *Widget) Shutdown() {
}

func (w *Widget) Params() any {
	return &w.params
}

func (w *Widget) Name() string {
	return name
}

func (w *Widget) format(d Data) (string, error) {
	return w.formatter.Format(w.params.Format, d)
}

func init() {
	widgets.Register(name, func() widgets.Widget {
		return &Widget{params: Params{
			Format:    `{{.Used | As "GiB" | Scale 1}} / {{.Total | As "GiB" | Scale 1}}`,
			LevelBad:  80,
			LevelGood: 50,
		}}
	})
}
