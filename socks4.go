package socks

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net"
	"strconv"
)

func (s *Server) handleSocks4(conn net.Conn) error {
	ctx := context.Background()
	var cddstportdstip [1 + 2 + 4]byte
	if _, err := io.ReadFull(conn, cddstportdstip[:]); err != nil {
		return err
	}
	command := cddstportdstip[0]
	dstPort := binary.BigEndian.Uint16(cddstportdstip[1:3])
	var dstIp net.IP = cddstportdstip[3:]
	if command != ConnectCommand {
		return fmt.Errorf("command %d is not supported", command)
	}
	// Skip USERID
	if _, err := readAsString(conn); err != nil {
		return err
	}
	// SOCKS4a
	if dstIp[0] == 0 && dstIp[1] == 0 && dstIp[2] == 0 && dstIp[3] != 0 {
		dstHost, err := readAsString(conn)
		_, dstIp, err = s.config.Resolver.Resolve(ctx, dstHost)
		if err != nil {
			return err
		}
	}
	if _, err := conn.Write([]byte{0, 90, 0, 0, 0, 0, 0, 0}); err != nil {
		return err
	}
	dial := s.config.Dial
	if dial == nil {
		dial = func(ctx context.Context, net_, addr string) (net.Conn, error) {
			return net.Dial(net_, addr)
		}
	}
	dstConn, err := dial(ctx, "tcp", net.JoinHostPort(dstIp.String(), strconv.Itoa(int(dstPort))))
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

func readAsString(r io.Reader) (string, error) {
	var buff bytes.Buffer
	var b [1]byte
	for {
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return "", err
		}
		if b[0] == 0 {
			break
		}
		buff.Write(b[:])
	}
	return buff.String(), nil
}
