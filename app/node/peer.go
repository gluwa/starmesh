package node

import (
	"fmt"
	"net"
)

// ConnectToPeer connects to another peer and sends a message.
func ConnectToPeer(address string, message string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte(message))

	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	fmt.Printf("[Client] Server response: %s\n", string(buffer[:n]))
}
