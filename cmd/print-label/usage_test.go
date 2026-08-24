package main

import (
	"strings"
	"testing"
)

func TestMetadataOutput(t *testing.T) {
	originalVersion := version
	version = "0.2.0"
	t.Cleanup(func() { version = originalVersion })

	usage, ok := metadataOutput([]string{"--usage"})
	if !ok {
		t.Fatal("--usage was not handled")
	}
	if !strings.Contains(usage, "version \"0.2.0\"") {
		t.Errorf("usage spec does not contain the version: %q", usage)
	}

	alias, ok := metadataOutput([]string{"--usage-spec"})
	if !ok || alias != usage {
		t.Error("--usage-spec did not return the same spec as --usage")
	}

	gotVersion, ok := metadataOutput([]string{"--version"})
	if !ok || gotVersion != "0.2.0\n" {
		t.Errorf("--version = %q, %v; want %q, true", gotVersion, ok, "0.2.0\n")
	}
}

func TestMetadataOutputIgnoresNormalArguments(t *testing.T) {
	for _, args := range [][]string{{}, {"label"}, {"--usage", "label"}} {
		if output, ok := metadataOutput(args); ok {
			t.Errorf("metadataOutput(%q) = %q, true; want unhandled", args, output)
		}
	}
}
