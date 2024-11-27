package main

import (
	"fmt"
	"net"
)

func main() {
	// pc, err := net.ListenPacket("udp", "localhost:40000")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// defer pc.Close()
	//
	// buffer := make([]byte, 100)
	// n, clientAddr, err := pc.ReadFrom(buffer)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	//
	// fmt.Println("Received", n, "bytes:", string(buffer), "from", clientAddr.String())

	l, err := net.Listen("tcp", "localhost:40000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Got connection from:", conn.RemoteAddr().String())

	var buf []byte
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
