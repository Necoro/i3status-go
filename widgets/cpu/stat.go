package cpu

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// The procfs package uses float64 to parse uint64 data -- which results in loss of precision.
// Thus, we have to implement it ourselves.

type Stat struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	Iowait  uint64
	Irq     uint64
	SoftIrq uint64
}

func (s Stat) Total() uint64 {
	return s.Idling() + s.Busy()
}

func (s Stat) Idling() uint64 {
	return s.Idle + s.Iowait
}
func (s Stat) Busy() uint64 {
	return s.User + s.Nice + s.System + s.Irq + s.SoftIrq
}

func safeSub(a, b uint64) uint64 {
	diff := a - b
	if diff > a { // underflow
		return 0
	}
	return diff
}

func (s Stat) Utilization(prev Stat) float64 {
	elapsed := safeSub(s.Total(), prev.Total())
	if elapsed == 0 {
		return 0
	}

	elapsedBusy := safeSub(s.Busy(), prev.Busy())

	return float64(elapsedBusy) / float64(elapsed)
}

func CollectStats() (Stat, []Stat, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return Stat{}, nil, fmt.Errorf("failed to open /proc/stat: %w", err)
	}
	defer f.Close()

	b := strings.Builder{}
	_, err = io.Copy(&b, f)
	if err != nil {
		return Stat{}, nil, fmt.Errorf("failed to read /proc/stat: %w", err)
	}

	globalStat := Stat{}
	stats := make([]Stat, 0)
	for line := range strings.Lines(b.String()) {
		if !strings.HasPrefix(line, "cpu") {
			continue
		}

		var cpu string
		stat := Stat{}
		_, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d", // ignore last 3 fields
			&cpu,
			&stat.User, &stat.Nice, &stat.System, &stat.Idle,
			&stat.Iowait, &stat.Irq, &stat.SoftIrq)

		if err != nil && err != io.EOF {
			return Stat{}, nil, fmt.Errorf("couldn't parse %q: %w", line, err)
		}

		if cpu == "cpu" {
			globalStat = stat
		} else {
			stats = append(stats, stat)
		}
	}

	return globalStat, stats, nil
}
