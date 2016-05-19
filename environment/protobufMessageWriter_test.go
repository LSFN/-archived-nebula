package environment

import (
	"bufio"
	"encoding/binary"
	"io"
	"testing"
)

func TestWriteMessages(t *testing.T) {
	pmw := NewProtobufMessageWriter(1)
	pipeReader, pipeWriter := io.Pipe()
	bufferChannel := make(chan []byte)
	go pmw.writeMessagesUntilClose(pipeWriter, bufferChannel)

	message := "Hello, World!"
	rawMessageLength := make([]byte, 8)
	lengthBytes := binary.PutUvarint(rawMessageLength, uint64(len(message)))
	rawMessageLength = rawMessageLength[:lengthBytes]
	rawMessage := append(rawMessageLength, []byte(message)...)
	bufferChannel <- rawMessage

	msgLength, err := binary.ReadUvarint(reader)
	if err != nil {
		// If there's an error, assume we aren't receiving the message and close up shop
		return
	}

	// If the message is of zero length, move straight to reading the next message
	if msgLength == 0 {
		continue
	}

	// Create a buffer to contain the message
	messageBuffer := make([]byte, msgLength)

	// Read the message proper
	if _, err := io.ReadFull(reader, messageBuffer); err != nil {
		return
	}

}
