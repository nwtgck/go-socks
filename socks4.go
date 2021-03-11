package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func handleSocks4(conn net.Conn) error {
	cddstportdstip := make([]byte, 1+2+4)
	if _, err := io.ReadFull(conn, cddstportdstip); err != nil {
		return err
	}
	command := cddstportdstip[0]
	dstPort := binary.BigEndian.Uint16(cddstportdstip[1:3])
	var dstIp net.IP = cddstportdstip[3:]
	if command != ConnectCommand {
		return fmt.Errorf("command %d is not supported", command)
	}
	// Skip USERID
	b := make([]byte, 1)
	for {
		if _, err := io.ReadFull(conn, b); err != nil {
			return err
		}
		if b[0] == 0 {
			break
		}
	}
	if _, err := conn.Write([]byte{0, 90, 0, 0, 0, 0, 0, 0}); err != nil {
		return err
	}
	dstConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dstIp, dstPort))
	if err != nil {
		return err
	}
	var errCh = make(chan error, 2)
	go func() {
		_, err := io.Copy(dstConn, conn)
		errCh <- err
	}()
	go func() {
		_, err := io.Copy(conn, dstConn)
		errCh <- err
	}()
	err = <-errCh
	if err != nil {
		return err
	}
	return <-errCh
}
