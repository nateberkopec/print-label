package main

import (
	"flag"
	"fmt"
	"os"

	printlabel "github.com/nateberkopec/print-label"
)

func main() {
	cfg, err := printlabel.LoadConfig(printlabel.ConfigPath)
	if err != nil {
		exit(err)
	}

	dryRun := flag.Bool("dry-run", false, "save a PNG preview instead of printing")
	printerIP := flag.String("printer-ip", cfg.PrinterIP, "printer IP address")
	printerPort := flag.Int("printer-port", cfg.PrinterPort, "printer TCP port")
	tapeWidth := flag.Float64("tape-width", cfg.TapeWidth, "tape width in mm")
	fontPath := flag.String("font", cfg.Font, "font path")
	fontIndex := flag.Int("font-index", cfg.FontIndex, "face index within font collections")
	fontSize := flag.String("font-size", cfg.FontSize, "font size in points, or auto")
	margin := flag.Int("margin", cfg.Margin, "left/right margin in pixels")
	out := flag.String("out", "label_preview.png", "preview path for --dry-run")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: print-label [options] text")
		flag.PrintDefaults()
		os.Exit(2)
	}

	cfg.PrinterIP = *printerIP
	cfg.PrinterPort = *printerPort
	cfg.TapeWidth = *tapeWidth
	cfg.Font = *fontPath
	cfg.FontIndex = *fontIndex
	cfg.FontSize = *fontSize
	cfg.Margin = *margin

	info, ok := printlabel.TapeInfos[cfg.TapeWidth]
	if !ok {
		exit(fmt.Errorf("unsupported tape width: %gmm", cfg.TapeWidth))
	}

	fmt.Printf("Tape: %gmm (%dpx height), margin: %dpx\n", cfg.TapeWidth, info.PrintablePixels, cfg.Margin)
	img, err := printlabel.RenderText(flag.Arg(0), cfg)
	if err != nil {
		exit(err)
	}
	fmt.Printf("Label image: %dx%dpx\n", img.Bounds().Dx(), img.Bounds().Dy())

	if *dryRun {
		if err := printlabel.SavePreview(img, *out, 3); err != nil {
			exit(err)
		}
		fmt.Printf("Preview saved to %s\n", *out)
		return
	}

	raster, err := printlabel.ImageToRaster(img, cfg.TapeWidth)
	if err != nil {
		exit(err)
	}
	job, err := printlabel.BuildPrintJob(raster, cfg.TapeWidth, cfg.Margin)
	if err != nil {
		exit(err)
	}
	fmt.Printf("Sending %d raster lines (%d bytes)...\n", len(raster), len(job))
	if err := printlabel.SendToPrinter(job, cfg.PrinterIP, cfg.PrinterPort); err != nil {
		exit(err)
	}
	fmt.Println("Print job sent.")
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
