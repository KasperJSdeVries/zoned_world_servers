package main

import (
	"fmt"
	"os"

	"github.com/maurice2k/tcpserver"

	"github.com/KasperJSdeVries/zoned_world_servers/internal/packet"
)

const Port = 40000

func main() {
	server, err := tcpserver.NewServer(fmt.Sprintf("localhost:%d", Port))
	if err != nil {
		fmt.Printf("error: could not start tcp server: %v\n", err)
		os.Exit(1)
	}

	server.SetRequestHandler(requestHandler)
	server.Listen()
	server.Serve()
}

func requestHandler(conn tcpserver.Connection) {
	for {
		packetData := packet.ReadPacket(conn)
		switch packetData[0] {
		}
	}
}
