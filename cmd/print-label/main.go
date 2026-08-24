package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	printlabel "github.com/nateberkopec/print-label"
)

func main() {
	if output, ok := metadataOutput(os.Args[1:]); ok {
		fmt.Print(output)
		return
	}

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

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: print-label [options] text [text ...]")
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
	images := make([]*image.Gray, flag.NArg())
	for i, text := range flag.Args() {
		images[i], err = printlabel.RenderText(text, cfg)
		if err != nil {
			exit(err)
		}
		fmt.Printf("Label %d image: %dx%dpx\n", i+1, images[i].Bounds().Dx(), images[i].Bounds().Dy())
	}

	if *dryRun {
		for i, img := range images {
			path := previewPath(*out, i, len(images))
			if err := printlabel.SavePreview(img, path, 3); err != nil {
				exit(err)
			}
			fmt.Printf("Preview saved to %s\n", path)
		}
		return
	}

	pages := make([][][]byte, len(images))
	totalLines := 0
	for i, img := range images {
		pages[i], err = printlabel.ImageToRaster(img, cfg.TapeWidth)
		if err != nil {
			exit(err)
		}
		totalLines += len(pages[i])
	}

	var job []byte
	if len(pages) == 1 {
		job, err = printlabel.BuildPrintJob(pages[0], cfg.TapeWidth, cfg.Margin)
	} else {
		job, err = printlabel.BuildMultiPrintJob(pages, cfg.TapeWidth, cfg.Margin)
	}
	if err != nil {
		exit(err)
	}
	fmt.Printf("Sending %d labels, %d raster lines (%d bytes)...\n", len(pages), totalLines, len(job))
	if err := printlabel.SendToPrinter(job, cfg.PrinterIP, cfg.PrinterPort); err != nil {
		exit(err)
	}
	fmt.Println("Print data sent; printer completion was not confirmed.")
}

func previewPath(path string, index, count int) string {
	if count == 1 {
		return path
	}
	ext := filepath.Ext(path)
	return fmt.Sprintf("%s-%d%s", strings.TrimSuffix(path, ext), index+1, ext)
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
