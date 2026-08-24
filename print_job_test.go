package printlabel

import (
	"bytes"
	"testing"
)

func TestUncompressedRasterTransfersEveryRow(t *testing.T) {
	lines := [][]byte{
		make([]byte, BytesPerLine),
		bytes.Repeat([]byte{0xff}, BytesPerLine),
	}
	var got bytes.Buffer
	writeRaster(&got, lines)

	var want bytes.Buffer
	for _, line := range lines {
		want.Write([]byte{0x47, BytesPerLine, 0x00})
		want.Write(line)
	}
	if !bytes.Equal(got.Bytes(), want.Bytes()) {
		t.Fatalf("uncompressed raster stream = %x, want %x", got.Bytes(), want.Bytes())
	}
}

func TestPrintJobMarksOnlyPageAsFinal(t *testing.T) {
	line := bytes.Repeat([]byte{0xff}, BytesPerLine)
	job, err := BuildPrintJob([][]byte{line}, 6, 14)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(job, []byte{0x1b, 0x69, 0x7a, 0xc4, 0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00}) {
		t.Fatal("single label is not marked as the final page")
	}
}

func TestMultiPrintJobCutsEachLabelWithoutLeaderFeed(t *testing.T) {
	line := bytes.Repeat([]byte{0xff}, BytesPerLine)
	job, err := BuildMultiPrintJob([][][]byte{{line}, {line}, {line}}, 6, 14)
	if err != nil {
		t.Fatal(err)
	}

	if got := bytes.Count(job, bytes.Repeat([]byte{0}, 200)); got != 1 {
		t.Fatalf("initialization count = %d, want 1", got)
	}
	if got := bytes.Count(job, []byte{0x1b, 0x69, 0x4d, 0x40}); got != 3 {
		t.Fatalf("auto-cut command count = %d, want 3", got)
	}
	if got := bytes.Count(job, []byte{0x1b, 0x69, 0x41, 0x01}); got != 3 {
		t.Fatalf("cut-each-label command count = %d, want 3", got)
	}
	if got := bytes.Count(job, []byte{0x1b, 0x69, 0x4b, 0x08}); got != 3 {
		t.Fatalf("full-cut command count = %d, want 3", got)
	}
	if got := bytes.Count(job, []byte{0x0c, 0x1b, 0x69, 0x61}); got != 2 {
		t.Fatalf("print-without-feed command count = %d, want 2", got)
	}
	if !bytes.Contains(job, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1b, 0x69, 0x4d}) {
		t.Fatal("first label is not marked as the starting page")
	}
	if !bytes.Contains(job, []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x1b, 0x69, 0x4d}) {
		t.Fatal("second label is not marked as an intermediate page")
	}
	if !bytes.Contains(job, []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x1b, 0x69, 0x4d}) {
		t.Fatal("third label is not marked as the final page")
	}
	if job[len(job)-1] != 0x1a {
		t.Fatalf("final command = %#x, want print-and-feed", job[len(job)-1])
	}
}
