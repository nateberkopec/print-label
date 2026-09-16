package printlabel

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"
)

func TestSendToPrinterWithoutStatusReply(t *testing.T) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	type result struct {
		data []byte
		err  error
	}
	results := make(chan result, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			results <- result{err: err}
			return
		}
		defer conn.Close()
		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			results <- result{err: err}
			return
		}
		data, err := io.ReadAll(conn)
		results <- result{data, err}
	}()

	job := bytes.Repeat([]byte{0x67, 0x46, 0x00}, 11000)
	addr := listener.Addr().(*net.TCPAddr)
	if err := SendToPrinter(job, addr.IP.String(), addr.Port); err != nil {
		t.Fatal(err)
	}
	received := <-results
	if received.err != nil {
		t.Fatal(received.err)
	}
	if !bytes.Equal(received.data, job) {
		t.Fatalf("received %d bytes, want the exact %d-byte print job", len(received.data), len(job))
	}
}
