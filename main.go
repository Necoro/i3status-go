package main

import (
	"fmt"

	"github.com/go-ini/ini"
)

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

	for _, b := range blocks {
		println(b.Run().Text)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
