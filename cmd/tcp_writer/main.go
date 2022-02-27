package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("hello, world")

	addr := net.TCPAddr{
		Port: 9000,
		IP:   net.ParseIP("0.0.0.0"),
	}
	listener, err := net.ListenTCP("tcp", &addr) // code does not block here
	if err != nil {
		panic(err)
	}
	conn, err := listener.AcceptTCP()
	if err != nil {
		panic(err)
	}
	file, err := os.Create("out.ts")
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	b := make([]byte, 2048)
	for i := 0; i < 1000; i++ {
		n, err := conn.Read(b)
		if err != nil {
			panic(err)
		}

		if n > 0 {
			write, err := file.Write(b[:n])
			if err != nil || write != n {
				panic(err)
			}
		}

	}

}
