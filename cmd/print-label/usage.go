package main

import "fmt"

var version = "dev"

func metadataOutput(args []string) (string, bool) {
	if len(args) != 1 {
		return "", false
	}

	switch args[0] {
	case "--usage", "--usage-spec":
		return usageSpec(), true
	case "--version":
		return version + "\n", true
	default:
		return "", false
	}
}

func usageSpec() string {
	return fmt.Sprintf(`min_usage_version "4.1.0"
name "print-label"
bin "print-label"
version %q
about "Print labels on a Brother PT-P900W using its raster TCP protocol"
repository "https://github.com/nateberkopec/print-label"

flag "--dry-run" help="Save a PNG preview instead of printing"
flag "--printer-ip <ip>" help="Printer IP address"
flag "--printer-port <port>" help="Printer TCP port"
flag "--tape-width <width>" help="Tape width in mm"
flag "--font <path>" help="Font path"
flag "--font-index <index>" help="Face index within font collections"
flag "--font-size <size>" help="Font size in points, or auto"
flag "--margin <pixels>" help="Left/right margin in pixels"
flag "--out <path>" help="Preview path for --dry-run"
flag "-h --help" help="Print help"
flag "--version" help="Print version"
flag "--usage --usage-spec" help="Print the CLI's usage specification"
arg "<text>..." help="Text for each label"
`, version)
}
