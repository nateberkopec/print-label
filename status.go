package printlabel

import "fmt"

const StatusFrameSize = 32

const (
	statusRequestReply      = 0x00
	statusPrintingCompleted = 0x01
	statusErrorOccurred     = 0x02
)

type Status struct {
	ExtendedError     byte
	ErrorInformation1 byte
	ErrorInformation2 byte
	MediaWidth        byte
	MediaType         byte
	StatusType        byte
	PhaseType         byte
	Notification      byte
}

func ParseStatus(frame []byte) (Status, error) {
	if len(frame) != StatusFrameSize {
		return Status{}, fmt.Errorf("invalid printer status length: got %d bytes, want %d", len(frame), StatusFrameSize)
	}
	if frame[0] != 0x80 || frame[1] != 0x20 || frame[2] != 0x42 || frame[3] != 0x30 {
		return Status{}, fmt.Errorf("invalid printer status header: %x", frame[:4])
	}
	return Status{
		ExtendedError:     frame[7],
		ErrorInformation1: frame[8],
		ErrorInformation2: frame[9],
		MediaWidth:        frame[10],
		MediaType:         frame[11],
		StatusType:        frame[18],
		PhaseType:         frame[19],
		Notification:      frame[22],
	}, nil
}

func (s Status) PrintingCompleted() bool {
	return s.StatusType == statusPrintingCompleted
}

func (s Status) ValidateTapeWidth(tapeWidth float64) error {
	info, ok := TapeInfos[tapeWidth]
	if !ok {
		return fmt.Errorf("unsupported tape width: %gmm", tapeWidth)
	}
	if s.MediaWidth != info.MediaByte {
		return fmt.Errorf("wrong media: requested %g mm tape, printer reports %g mm", tapeWidth, statusMediaWidth(s.MediaWidth))
	}
	return nil
}

func statusMediaWidth(value byte) float64 {
	if value == 0x04 {
		return 3.5
	}
	return float64(value)
}
