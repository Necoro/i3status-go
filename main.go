package main

import (
	"fmt"

	"github.com/go-ini/ini"
)

func run() error {
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

	for _, section := range cfg.Sections() {
		_, _ = NewBlock(section, defaultOptions)
		println(section.Name())
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
