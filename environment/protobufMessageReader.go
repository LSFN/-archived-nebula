// protobufMessageReader
package environment

import (
	"bufio"
	"encoding/binary"
	"io"
)

type protobufMessageReader struct {
	maxBufferedMessages uint
}

func NewProtobufMessageReader(maxBufferedMessages uint) *protobufMessageReader {
	msgReader := new(protobufMessageReader)
	msgReader.maxBufferedMessages = maxBufferedMessages
	return msgReader
}

func (p *protobufMessageReader) Start(reader io.Reader) <-chan []byte {
	bufferChannel := make(chan []byte, p.maxBufferedMessages)
	go p.readMessagesUntilError(bufio.NewReader(reader), bufferChannel)
	return bufferChannel
}

func (p *protobufMessageReader) readMessagesUntilError(reader *bufio.Reader, bufferChannel chan<- []byte) {
	// When we leave here, close the channel to indicate we are done.
	defer close(bufferChannel)

	for {
		// Read the length of the protobuf message to receive from the reader
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

		bufferChannel <- messageBuffer
	}
}
