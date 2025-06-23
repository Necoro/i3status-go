package config

import (
	"fmt"
	"os"

	"github.com/Necoro/i3status-go/config/parser"
)

type Config = parser.Config
type Params = parser.Params
type Section = parser.Section

func Load(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file '%s': %w", filename, err)
	}

	cfg, err := parser.ParseReader(filename, f)
	if err != nil {
		return nil, err
	}

	return cfg.(*Config), nil
}
