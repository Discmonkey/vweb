package utils

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func HttpNotOk(statusCode int, w http.ResponseWriter, frontendErr string, err error) bool {
	if err != nil {
		http.Error(w, frontendErr, statusCode)
		log.Println(err)
		return true
	} else {
		return false
	}
}

// NewRandomUdpConn returns a fresh UDP connection and the port of said connection
func NewRandomUdpConn() (net.Conn, int, error) {
	address := net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 0,
	}
	udpConn, err := net.ListenUDP("udp", &address)
	if err != nil {
		return nil, 0, err
	}

	port, err := Port(udpConn.LocalAddr().String())
	if err != nil {
		return nil, 0, err
	}
	return udpConn, port, nil
}

func Port(address string) (int, error) {
	index := strings.LastIndex(address, ":")
	return strconv.Atoi(address[index+1:])
}
