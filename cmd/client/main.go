package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:40000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()


}
