package printlabel

import (
	"bytes"
	"io"
	"net"
	"strconv"
	"time"
)

func SendToPrinter(data []byte, ip string, port int) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(port)), 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	return err
}
