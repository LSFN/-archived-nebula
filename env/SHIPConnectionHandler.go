// SHIPConnectionHandler
package env

import (
	"fmt"
	"net"
)

type SHIPConnectionHandlerInputMessageType int

const (
	SCH_DISCONNECT = iota
)

type SHIPConnectionHandlerInputMessage struct {
	messageType SHIPConnectionHandlerInputMessageType
}

type SHIPConnectionHandlerOutputMessageType int

const (
	SCH_DISCONNECTED = iota
)

type SHIPConnectionHandlerOutputMessage struct {
	messageType SHIPConnectionHandlerOutputMessageType
}

type SHIPConnectionHandler struct {
	inputMessageChannel  chan SHIPConnectionHandlerInputMessage
	outputMessageChannel chan SHIPConnectionHandlerOutputMessage
}

func (c *SHIPConnectionHandler) Start(conn net.Conn) {
	c.inputMessageChannel = make(chan SHIPConnectionHandlerInputMessage)
	c.outputMessageChannel = make(chan SHIPConnectionHandlerOutputMessage)
	buf := make([]byte, 0, 1024)
	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection read error.")
			conn.Close()
			break
		}
		fmt.Printf("Read %d bytes", bytesRead)
	}
}
