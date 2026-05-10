package printlabel

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/font/gofont/gobold"
)

func testConfig(t *testing.T) Config {
	t.Helper()
	path := filepath.Join(t.TempDir(), "font.ttf")
	if err := os.WriteFile(path, gobold.TTF, 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := DefaultConfig()
	cfg.Font = path
	cfg.FontIndex = 0
	return cfg
}
