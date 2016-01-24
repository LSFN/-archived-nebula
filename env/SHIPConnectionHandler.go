// SHIPConnectionHandler
package env

import (
	"net"
	"fmt"
)

type SHIPMessage int

type SHIPConnectionHandler struct {
	connection net.Conn
	IncomingMessages chan SHIPMessage
}

func (c *SHIPConnectionHandler) HandleSHIPConnection(conn net.Conn) {
	c.connection = conn
	for {
		buf := make([]byte, 0, 1024)
		bytesRead, err := c.connection.Read(buf)
		if err != nil {
			fmt.Println("Connection read error.")
		}
		fmt.Printf("Read %d bytes", bytesRead)
	}
}
