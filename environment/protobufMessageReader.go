// protobufMessageReader
package environment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MINIMUM_READ_MESSAGE_SIZE = 255
)

type protobufMessageReader struct {
	maxBytesPerMessage  uint64
	bytesForLength      uint8
	maxBufferedMessages uint
}

func NewProtobufMessageReader(maxBytesPerMessage uint64, maxBufferedMessages uint) *protobufMessageReader {
	msgReader := new(protobufMessageReader)
	msgReader.maxBytesPerMessage = maxBytesPerMessage
	msgReader.maxBufferedMessages = maxBufferedMessages

	// Constrain maxBufferedMessages to a minimum
	if msgReader.maxBufferedMessages < 1 {
		msgReader.maxBufferedMessages = 1
	}

	// Constrain maxBytesPerMessage to a minimum
	if msgReader.maxBytesPerMessage < MINIMUM_READ_MESSAGE_SIZE {
		msgReader.maxBytesPerMessage = MINIMUM_READ_MESSAGE_SIZE
	}

	// Figure out the number of bytes needed to represent that maximum size
	msgReader.bytesForLength = 1
	var maxLengthRepresentable uint64 = 255
	for maxLengthRepresentable < maxBytesPerMessage && msgReader.bytesForLength < 8 {
		msgReader.bytesForLength++
		maxLengthRepresentable <<= 8
		maxLengthRepresentable += 255
	}
	// As nothing greater than MaxUint64 can be passed in, the number of bytes required will never exceed 8.
	// Whether or not the host system has enough memory for such a buffer is someone else's problem

	return msgReader
}

func (p *protobufMessageReader) Start(reader io.Reader) <-chan []byte {
	bufferChannel := make(chan []byte, p.maxBufferedMessages)
	go p.readMessagesUntilError(reader, bufferChannel)
	return bufferChannel
}

func (p *protobufMessageReader) readMessagesUntilError(reader io.Reader, bufferChannel chan<- []byte) {
	// When we leave here, close the channel to indicate we are done.
	defer close(bufferChannel)

	lengthBuffer := make([]byte, p.bytesForLength)

	for {
		// Read the length of the protobuf message to receive from the reader
		fmt.Println("Reading length")
		var readLengthBytes uint8 = 0
		for readLengthBytes < p.bytesForLength {
			n, err := reader.Read(lengthBuffer[readLengthBytes:])
			readLengthBytes += uint8(n)

			// If there's an error, assume we aren't receiving the message and close up shop
			if err != nil {
				return
			}
		}

		// Convert the length to a number
		fmt.Println("Converting length")
		lengthReader := bytes.NewReader(lengthBuffer)
		var msgLength uint64
		if err := binary.Read(lengthReader, binary.LittleEndian, &msgLength); err != nil {
			fmt.Println("error", err)
		}
		fmt.Printf("Length is %d\n", msgLength)

		// If the message is of zero length, move straight to reading the next message
		if msgLength == 0 {
			fmt.Println("Message length was 0, continuing")
			continue
		}

		// If the received length is greater than our max length, cop out
		if msgLength > p.maxBytesPerMessage {
			fmt.Println("Message length was above maximum of" p.maxBytesPerMessage)
			return
		}

		// Create a buffer to contain the message
		messageBuffer := make([]byte, msgLength)

		// Read the message proper
		fmt.Println("Reading message")
		var readMsgBytes uint64 = 0
		for readMsgBytes < msgLength {
			n, err := reader.Read(messageBuffer[readMsgBytes:])
			readMsgBytes += uint64(n)

			// If the message has been received completely,
			// pass it along regardless of error
			if readMsgBytes == msgLength {
				fmt.Println("Read complete message")
				bufferChannel <- messageBuffer
			}

			// If there's an error
			if err != nil {
				return
			}
		}
	}
}
