package cpu

import (
	"fmt"
	"math"
	"slices"

	u "github.com/Necoro/go-units"
	"github.com/prometheus/procfs"

	"github.com/Necoro/i3status-go/widgets"
)

const (
	name      = "cpu"
	icon rune = '\uf085' // 
)

type Params struct {
	// Format the output string. Takes a Go template string.
	// For possible values to be used, see [Data].
	Format string
	// LevelGood is the upper bound for the good level.
	LevelGood int
	// LevelBad is the lower bound for the bad level.
	LevelBad int
}

type Widget struct {
	params         Params
	formatter      widgets.Formatter[Data]
	barchart       barchart
	prevCpuStats   []Stat
	prevAvgCpuStat Stat
}

type Data struct {
	Frequencies    []u.Value
	Utilization    []u.Value
	AvgUtilization u.Value
}

func (d Data) UtilChartAvg() string {
	return string(barchars.Bar(d.AvgUtilization.Float()))
}

func (d Data) UtilChart() string {
	bars := make([]rune, len(d.Utilization))

	for i, u := range d.Utilization {
		bars[i] = barchars.Bar(u.Float())
	}

	return string(bars)
}

func valueAvg(vals []u.Value) float64 {
	sum := 0.
	for _, v := range vals {
		sum += v.Float()
	}
	return sum / float64(len(vals))
}

func (d Data) UtilChartCombined(combine int) string {
	l := int(
		math.Ceil(
			float64(len(d.Utilization)) / float64(combine)))

	bars := make([]rune, l)

	i := 0
	for us := range slices.Chunk(d.Utilization, combine) {
		bars[i] = barchars.Bar(valueAvg(us))
		i++
	}

	return string(bars)
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

	avgCpuStats, cpuStats, err := CollectStats()
	if err != nil {
		return
	}

	if len(cpuStats) != len(cpuinfo) {
		err = fmt.Errorf("cpuStats and cpuinfo have different length")
		return
	}

	if w.prevCpuStats == nil {
		w.prevCpuStats = cpuStats
		w.prevAvgCpuStat = avgCpuStats
	}

	data := Data{
		Frequencies:    make([]u.Value, len(cpuinfo)),
		Utilization:    make([]u.Value, len(cpuinfo)),
		AvgUtilization: u.NewValue(avgCpuStats.Utilization(w.prevAvgCpuStat)*100, u.Percent),
	}
	for i := 0; i < len(cpuinfo); i++ {
		data.Frequencies[i] = u.NewValue(cpuinfo[i].CPUMHz, u.MegaHertz)
		data.Utilization[i] = u.NewValue(cpuStats[i].Utilization(w.prevCpuStats[i])*100, u.Percent)
	}

	w.prevCpuStats = cpuStats
	w.prevAvgCpuStat = avgCpuStats

	avgUtil := int(data.AvgUtilization.Float())
	if avgUtil < w.params.LevelGood {
		d.State = widgets.StateGood
	} else if avgUtil >= w.params.LevelBad {
		d.State = widgets.StateBad
	} else {
		d.State = widgets.StateMid
	}

	d.Text, err = w.format(data)
	return
}

func (w *Widget) format(d Data) (string, error) {
	return w.formatter.Format(w.params.Format, d)
}

func (w *Widget) Shutdown() {
}

func (w *Widget) Params() any {
	return &w.params
}

func init() {
	widgets.Register(name, func() widgets.Widget {
		return &Widget{
			params: Params{
				Format: "{{.UtilChartAvg}} {{index .Frequencies 0}} {{index .Frequencies 1}}",
			}}
	})
}
