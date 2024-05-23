package utils

import "net"

func ReadFromConnection(buffer []byte, conn net.Conn) ([]byte, error) {
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}
