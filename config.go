package printlabel

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	PrinterIP   string
	PrinterPort int
	TapeWidth   float64
	Font        string
	FontIndex   int
	FontSize    string
	Margin      int
}

func DefaultConfig() Config {
	return Config{
		PrinterIP:   "192.168.1.86",
		PrinterPort: 9100,
		TapeWidth:   12,
		Font:        "/System/Library/Fonts/Helvetica.ttc",
		FontIndex:   1,
		FontSize:    "auto",
		Margin:      14,
	}
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	f, err := os.Open(ExpandHome(path))
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		if err := applyConfigLine(&cfg, s.Text()); err != nil {
			return cfg, err
		}
	}
	return cfg, s.Err()
}

func applyConfigLine(cfg *Config, raw string) error {
	line := strings.TrimSpace(stripComment(raw))
	if line == "" {
		return nil
	}
	key, value, ok := strings.Cut(line, ":")
	if !ok {
		return nil
	}
	return setConfigValue(cfg, strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), `"'`))
}

func setConfigValue(cfg *Config, key, value string) error {
	var err error
	switch key {
	case "printer_ip":
		cfg.PrinterIP = value
	case "printer_port":
		cfg.PrinterPort, err = strconv.Atoi(value)
	case "tape_width":
		cfg.TapeWidth, err = strconv.ParseFloat(value, 64)
	case "font":
		cfg.Font = value
	case "font_index":
		cfg.FontIndex, err = strconv.Atoi(value)
	case "font_size":
		cfg.FontSize = value
	case "margin":
		cfg.Margin, err = strconv.Atoi(value)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", key, err)
	}
	return nil
}

func ExpandHome(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	return path
}
