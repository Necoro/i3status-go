package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-ini/ini"
)

var i3Output = os.Stdout

func writeStatus(blocks []*Block) error {
	encoder := json.NewEncoder(i3Output)
	encoder.SetEscapeHTML(false) // may contain Pango markup
	encoder.SetIndent("", "")    // don't indent

	hdr := I3BarHeader{
		Version:     1,        // required
		ClickEvents: false,    // not supported atm
		StopSignal:  new(int), // ignore this for the moment
	}

	if err := encoder.Encode(hdr); err != nil {
		return fmt.Errorf("failed to encode header: %w", err)
	}

	if _, err := i3Output.WriteString("[[]\n,"); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	for {
		msgs := make([]I3BarBlock, len(blocks))
		for i, b := range blocks {
			d := b.Run()
			msgs[i] = NewI3BarBlock(d)
		}
		if err := encoder.Encode(msgs); err != nil {
			return fmt.Errorf("failed to encode blocks: %w", err)
		}

		if _, err := i3Output.WriteString(","); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}

		time.Sleep(5 * time.Second)
	}
}

func run() error {
	ini.DefaultSection = "default" // else it does not work with `Insensitive`

	options := ini.LoadOptions{
		AllowNonUniqueSections: true,
		Insensitive:            true,
	}

	cfg, err := ini.LoadSources(options, "test.config")
	if err != nil {
		return fmt.Errorf("failed to load config '%s': %w", "test.config", err)
	}

	defaultOptions, err := NewBlock(cfg.Section(ini.DefaultSection), DefaultBlock)
	if err != nil {
		return fmt.Errorf("failed to load default options: %w", err)
	}

	blocks := make([]*Block, 0, len(cfg.Sections())-1)
	for _, section := range cfg.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}

		b, err := NewBlock(section, defaultOptions)
		if err != nil {
			return fmt.Errorf("failed to load block '%s': %w", section.Name(), err)
		}

		blocks = append(blocks, b)
	}

	return writeStatus(blocks)
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
