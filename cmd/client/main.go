package main

import (
	"fmt"
	"net"
	"github.com/KasperJSdeVries/zoned_world_servers/internal/common/vec"
)

func main() {
	// raddr, err := net.ResolveUDPAddr("udp", "localhost:40000")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	//
	// conn, err := net.DialUDP("udp", nil, raddr)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// defer conn.Close()
	//
	// reader := strings.NewReader("test")
	// _, err = io.Copy(conn, reader)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	position := vec.Vec2 {}

	conn, err := net.Dial("tcp", "localhost:40000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	conn.Write()
}
