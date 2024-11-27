package packet

import (
	"encoding/binary"
	"io"
	"net"
)

func ReadPacket(conn net.Conn) []byte {
	var buffer [4]byte
	index := 0
	for {
		n, err := conn.Read(buffer[index:4])
		if err == io.EOF {
			return nil
		}
		index += n
		if index == 4 {
			break
		}
	}

	length := binary.LittleEndian.Uint32(buffer[:])

	if length == 0 {
		return nil
	}

	packetData := make([]byte, length)
	index = 0
	for {
		n, err := conn.Read(packetData[index:length])
		if err == io.EOF {
			return nil
		}
		index += n
		if index == int(length) {
			break
		}
	}

	return packetData
}
