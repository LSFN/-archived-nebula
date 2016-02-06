// SHIPConnectionHandler
package env

import (
	"bufio"
	"encoding/binary"
	"net"

	"github.com/golang/protobuf/proto"

	"github.com/LSFN/shipenvproto"
)

const (
	DEFAULT_READ_BUFFER_SIZE = 4096
)

type SHIPConnectionHandler struct {
	inboundMessages  chan *shipenvproto.SHIPtoENV
	outboundMessages chan *shipenvproto.ENVtoSHIP
}

func (c *SHIPConnectionHandler) Start(conn net.Conn) {
	c.inboundMessages = make(chan *shipenvproto.SHIPtoENV)
	c.outboundMessages = make(chan *shipenvproto.ENVtoSHIP)

	go c.readMessages(conn)
	go c.writeMessages(conn)
}

func (c *SHIPConnectionHandler) readMessages(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	readBuffer := make([]byte, DEFAULT_READ_BUFFER_SIZE)

readLoop:
	for {
		// Read the message length
		msgLen, err := binary.ReadUvarint(reader)
		if err != nil {
			break
		}

		// Expand main buffer if necessary
		currentCap := uint64(cap(readBuffer))
		if currentCap < msgLen {
			for currentCap < msgLen {
				currentCap *= 2
			}
			readBuffer = make([]byte, currentCap)
		}

		// Read the message body
		var bytesRead uint64 = 0
		msgReadBuffer := readBuffer[:msgLen]
		for bytesRead < msgLen {
			n, err := reader.Read(msgReadBuffer[bytesRead:])
			bytesRead += uint64(n)
			if err != nil {
				break readLoop
			}
		}

		// Unmarshal the message
		message := new(shipenvproto.SHIPtoENV)
		if err := proto.Unmarshal(msgReadBuffer, message); err != nil {
			break
		}

		// Send message on channel
		c.inboundMessages <- message
	}

	// The SHIP has disconnected or the connection has suffered an error
	// Close the inbound channel to indicate this to the next layer
	close(c.inboundMessages)
}

func (c *SHIPConnectionHandler) writeMessages(conn net.Conn) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)

writeLoop:
	// Loop whilst the channel is open
	for message := range c.outboundMessages {
		// Marshal the message
		marshaledMessage, err := proto.Marshal(message)
		if err != nil {
			break
		}

		// Write the length of the message
		lengthBuffer := make([]byte, 0, binary.MaxVarintLen64)
		binary.PutUvarint(lengthBuffer, uint64(len(marshaledMessage)))
		bytesWritten := 0
		for bytesWritten < len(lengthBuffer) {
			n, err := writer.Write(lengthBuffer)
			bytesWritten += n
			if err != nil {
				break writeLoop
			}
		}

		// Write the message body
		bytesWritten = 0
		for bytesWritten < len(marshaledMessage) {
			n, err := writer.Write(marshaledMessage)
			bytesWritten += n
			if err != nil {
				break writeLoop
			}
		}
	}
}
