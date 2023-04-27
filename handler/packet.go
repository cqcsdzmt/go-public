package handler

import (
	"fmt"
	"go-public/common"
	"io"
	"net"
)

var tokenSize = 16

var (
	HelloPacket = []byte("\x00")[0]
	ConnPacket  = []byte("\x01")[0]
)

var (
	HelloPacketSize = 1 + 2 + tokenSize
	ConnPacketSize  = 1 + tokenSize
)

func sendHelloPacket(conn net.Conn, remotePort int) error {
	port := uint16(remotePort)
	var tokenBytes = common.Token2Bytes(common.ClientConfig.Token)
	var buf = make([]byte, HelloPacketSize)
	buf[0] = HelloPacket
	buf[1] = byte(port >> 8)
	buf[2] = byte(port)
	copy(buf[3:], tokenBytes)
	_, err := conn.Write(buf)
	return err
}

func parseHelloPacket(conn net.Conn) (bool, int) {
	port := make([]byte, 2)
	n, err := io.ReadFull(conn, port)
	if n != 2 || err != nil {
		return false, 0
	}

	portNumber := int(port[0])<<8 + int(port[1])
	if portNumber == common.ServerConfig.Port {
		fmt.Println("Invalid port:", portNumber)
		return false, 0
	}

	token := make([]byte, tokenSize)
	n, err = io.ReadFull(conn, token)
	if n != tokenSize || err != nil {
		return false, 0
	}

	tokenString := common.Bytes2Token(token)
	if tokenString == common.ServerConfig.Token {
		return true, portNumber
	}
	fmt.Println("Invalid token:", tokenString)
	return false, 0
}

func SendConnPacket(conn net.Conn, token string) error {
	var buf = make([]byte, ConnPacketSize)
	buf[0] = ConnPacket
	copy(buf[1:], common.Token2Bytes(token))
	_, err := conn.Write(buf)
	return err
}
