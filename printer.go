package printlabel

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func SendToPrinter(data []byte, tapeWidth float64, ip string, port int) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := validatePrinterStatus(conn, tapeWidth); err != nil {
		return err
	}
	_, err = io.Copy(conn, bytes.NewReader(data))
	return err
}

func validatePrinterStatus(conn net.Conn, tapeWidth float64) error {
	if _, err := conn.Write([]byte{0x1b, 0x69, 0x53}); err != nil {
		return err
	}
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}

	frame := make([]byte, StatusFrameSize)
	if _, err := io.ReadFull(conn, frame); err != nil {
		return err
	}
	status, err := ParseStatus(frame)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	if status.StatusType != statusRequestReply {
		return fmt.Errorf("unexpected printer status type: 0x%02x", status.StatusType)
	}
	return status.ValidateTapeWidth(tapeWidth)
}
