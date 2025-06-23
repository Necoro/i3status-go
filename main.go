package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Necoro/i3status-go/config"
)

var i3Output = os.Stdout

func runBlocks(encoder *json.Encoder, blocks []*Block) error {
	msgs := make([]I3BarBlock, len(blocks))

	var wg sync.WaitGroup
	wg.Add(len(blocks))

	for i, b := range blocks {
		go func() {
			d := b.Run()
			msgs[i] = NewI3BarBlock(b, d)
			wg.Done()
		}()
	}

	wg.Wait()

	if err := encoder.Encode(msgs); err != nil {
		return fmt.Errorf("failed to encode blocks: %w", err)
	}

	if _, err := i3Output.WriteString(","); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

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

	// first
	if err := runBlocks(encoder, blocks); err != nil {
		return err
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for _ = range ticker.C {
		if err := runBlocks(encoder, blocks); err != nil {
			return err
		}
	}

	return nil
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("no config file given")
	}

	fileName := os.Args[1]
	cfg, err := config.Load(fileName)
	if err != nil {
		return fmt.Errorf("failed to load config '%s': %w", "test.config", err)
	}

	defaultOptions, err := GlobalBlock(cfg.Global)
	if err != nil {
		return fmt.Errorf("failed to load default options: %w", err)
	}

	blocks := make([]*Block, 0, len(cfg.Sections))
	for _, section := range cfg.Sections {
		b, err := NewBlock(section, defaultOptions)
		if err != nil {
			return fmt.Errorf("failed to load block '%s': %w", section.FullName(), err)
		}

		blocks = append(blocks, b)
	}

	if err = writeStatus(blocks); err != nil {
		return err
	}

	for _, b := range blocks {
		b.Widget.Shutdown()
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
