package printlabel

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"
)

func TestSendToPrinterChecksTapeBeforeSendingJob(t *testing.T) {
	printer, results := startTestPrinter(t, validStatusFrame(statusRequestReply))
	job := []byte("print job")

	if err := SendToPrinter(job, 6, printer.IP.String(), printer.Port); err != nil {
		t.Fatal(err)
	}
	result := <-results
	if result.err != nil {
		t.Fatal(result.err)
	}
	if !bytes.Equal(result.request, []byte{0x1b, 0x69, 0x53}) {
		t.Errorf("status request = %x, want 1b6953", result.request)
	}
	if !bytes.Equal(result.job, job) {
		t.Errorf("print job = %q, want %q", result.job, job)
	}
}

func TestSendToPrinterRejectsInvalidStatusBeforeSendingJob(t *testing.T) {
	tests := []struct {
		name  string
		frame []byte
		width float64
		want  string
	}{
		{"wrong tape", validStatusFrame(statusRequestReply), 12, "printer reports 6 mm"},
		{"printer error", statusFrameWithError(), 6, "cover open"},
		{"unexpected status", validStatusFrame(statusPrintingCompleted), 6, "unexpected printer status type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer, results := startTestPrinter(t, tt.frame)
			err := SendToPrinter([]byte("must not be sent"), tt.width, printer.IP.String(), printer.Port)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("SendToPrinter error = %v, want %q", err, tt.want)
			}
			result := <-results
			if result.err != nil {
				t.Fatal(result.err)
			}
			if len(result.job) != 0 {
				t.Errorf("sent print job despite invalid status: %q", result.job)
			}
		})
	}
}

func statusFrameWithError() []byte {
	frame := validStatusFrame(statusErrorOccurred)
	frame[9] = 0x10
	return frame
}

type printerResult struct {
	request []byte
	job     []byte
	err     error
}

func startTestPrinter(t *testing.T, status []byte) (*net.TCPAddr, <-chan printerResult) {
	t.Helper()
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatal(err)
	}

	results := make(chan printerResult, 1)
	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			results <- printerResult{err: err}
			return
		}
		defer conn.Close()

		request := make([]byte, 3)
		if _, err := io.ReadFull(conn, request); err != nil {
			results <- printerResult{err: err}
			return
		}
		if _, err := conn.Write(status); err != nil {
			results <- printerResult{request: request, err: err}
			return
		}
		job, err := io.ReadAll(conn)
		results <- printerResult{request: request, job: job, err: err}
	}()

	return listener.Addr().(*net.TCPAddr), results
}
