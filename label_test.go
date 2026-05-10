package printlabel

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`printer_ip: "10.0.0.9"
printer_port: 9101
tape_width: 18
font: "/tmp/Label Font.ttc"
font_index: 2 # bold
font_size: auto
margin: 22
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PrinterIP != "10.0.0.9" || cfg.PrinterPort != 9101 || cfg.TapeWidth != 18 || cfg.Font != "/tmp/Label Font.ttc" || cfg.FontIndex != 2 || cfg.FontSize != "auto" || cfg.Margin != 22 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}

func TestRenderTextGoldenPNGs(t *testing.T) {
	for _, text := range []string{"Birdie", "Sky"} {
		t.Run(text, func(t *testing.T) {
			img, err := RenderText(text, DefaultConfig())
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			if err := png.Encode(&buf, img); err != nil {
				t.Fatal(err)
			}
			assertGolden(t, filepath.Join("testdata", "golden", text+".png"), buf.Bytes())
		})
	}
}

func TestRasterAndPrintJobGolden(t *testing.T) {
	img, err := RenderText("Birdie", DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	raster, err := ImageToRaster(img, DefaultConfig().TapeWidth)
	if err != nil {
		t.Fatal(err)
	}
	job, err := BuildPrintJob(raster, DefaultConfig().TapeWidth, DefaultConfig().Margin)
	if err != nil {
		t.Fatal(err)
	}

	var rasterBytes []byte
	for _, line := range raster {
		rasterBytes = append(rasterBytes, line...)
	}
	assertGoldenText(t, filepath.Join("testdata", "golden", "Birdie.raster.sha256"), sha256Hex(rasterBytes)+"\n")
	assertGoldenText(t, filepath.Join("testdata", "golden", "Birdie.printjob.sha256"), sha256Hex(job)+"\n")
}

func assertGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s mismatch: got sha256=%s want sha256=%s; run go test ./... -update to accept", path, sha256Hex(got), sha256Hex(want))
	}
}

func assertGoldenText(t *testing.T, path string, got string) {
	t.Helper()
	assertGolden(t, path, []byte(got))
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
