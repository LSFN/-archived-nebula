// SHIPConnectionHandler
package env

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"

	"github.com/LSFN/shipenvproto"
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

const (
	DEFAULT_READ_BUFFER_SIZE = 4096
)

type SHIPConnectionHandler struct {
	inputMessageChannel  chan SHIPConnectionHandlerInputMessage
	outputMessageChannel chan SHIPConnectionHandlerOutputMessage
	inboundMessages      chan *shipenvproto.SHIPtoENV
	outboundMessages     chan *shipenvproto.ENVtoSHIP
}

func (c *SHIPConnectionHandler) Start(conn net.Conn) {
	defer conn.Close()
	c.inputMessageChannel = make(chan SHIPConnectionHandlerInputMessage)
	c.outputMessageChannel = make(chan SHIPConnectionHandlerOutputMessage)
	c.inboundMessages = make(chan *shipenvproto.SHIPtoENV)
	c.outboundMessages = make(chan *shipenvproto.ENVtoSHIP)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readBuffer := make([]byte, DEFAULT_READ_BUFFER_SIZE)

outer:
	for {
		// Read the message length
		msgLen, err := binary.ReadUvarint(reader)
		if err != nil {
			break
		}

		// Expand main buffer if necessary
		currentCap := cap(readBuffer)
		if currentCap < msgLen {
			for currentCap < msgLen {
				currentCap *= 2
			}
			readBuffer = make([]byte, currentCap)
		}

		// Read the message body
		bytesRead := 0
		msgReadBuffer := readBuffer[:msgLen]
		for bytesRead < msgLen {
			n, err := reader.Read(msgReadBuffer[bytesRead:])
			bytesRead += n
			if err != nil {
				break outer
			}
		}

		// Unmarshal the message
		message := new shipenvproto.ENVtoSHIP
		if err := proto.Unmarshal(msgReadBuffer, message); err != nil {
			break
		}
		
		// Send message on channel
		c.inboundMessages <- message
	}
	
}
