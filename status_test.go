package printlabel

import (
	"strings"
	"testing"
)

func TestStatusErrors(t *testing.T) {
	tests := []struct {
		name  string
		index int
		value byte
		want  string
	}{
		{"no media", 8, 0x01, "no media"},
		{"cutter jam", 8, 0x04, "cutter jam"},
		{"wrong media", 9, 0x01, "wrong media"},
		{"communication error", 9, 0x04, "communication error"},
		{"cover open", 9, 0x10, "cover open"},
		{"unknown extended error", 7, 0xff, "unknown extended error 0xff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := validStatusFrame(0x02)
			frame[tt.index] = tt.value
			status, err := ParseStatus(frame)
			if err != nil {
				t.Fatal(err)
			}
			if err := status.Err(); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("status error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestStatusValidatesLoadedTapeWidth(t *testing.T) {
	status, err := ParseStatus(validStatusFrame(0x00))
	if err != nil {
		t.Fatal(err)
	}
	if err := status.ValidateTapeWidth(12); err == nil || !strings.Contains(err.Error(), "6 mm") {
		t.Fatalf("tape width error = %v, want loaded 6 mm", err)
	}
}

func TestPrintingCompletedStatus(t *testing.T) {
	status, err := ParseStatus(validStatusFrame(0x01))
	if err != nil {
		t.Fatal(err)
	}
	if err := status.Err(); err != nil {
		t.Fatal(err)
	}
	if !status.PrintingCompleted() {
		t.Fatal("printing-completed frame was not recognized")
	}
}

func validStatusFrame(statusType byte) []byte {
	frame := make([]byte, StatusFrameSize)
	copy(frame, []byte{0x80, 0x20, 0x42, 0x30, 0x69, 0x30})
	frame[10] = 6
	frame[11] = 1
	frame[18] = statusType
	return frame
}
