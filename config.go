package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Necoro/i3status-go/parser"
)

type Params map[string]string

type Config struct {
	Global   Params
	Sections []Section
}

type Section struct {
	Name      string
	Qualifier string
	Params    Params
}

func (s Section) FullName() string {
	return s.Name + "." + s.Qualifier
}

func params(p []parser.Parameter) Params {
	params := make(Params, len(p))
	for _, param := range p {
		params[strings.ToLower(param.Name)] = param.Value
	}
	return params
}

func fromParsedConfig(cfg *parser.Config) *Config {
	c := Config{}
	c.Global = params(cfg.GlobalParams)
	c.Sections = make([]Section, len(cfg.Sections))
	for i, section := range cfg.Sections {
		c.Sections[i] = Section{
			Name:      strings.ToLower(section.Name),
			Qualifier: strings.ToLower(section.Qualifier),
			Params:    params(section.Params),
		}
	}
	return &c
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file '%s': %w", filename, err)
	}

	cfg, err := parser.Parse(filename, f)
	if err != nil {
		return nil, err
	}
	return fromParsedConfig(cfg), nil
}
