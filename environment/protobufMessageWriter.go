// protobufMessageWriter
package environment

import (
	"encoding/binary"
	"io"
)

type protobufMessageWriter struct {
	maxBufferedMessages uint
}

func NewProtobufMessageWriter(maxBufferedMessages uint) *protobufMessageWriter {
	msgWriter := new(protobufMessageWriter)
	msgWriter.maxBufferedMessages = maxBufferedMessages
	return msgWriter
}

func (p *protobufMessageWriter) Start(writer io.Writer) chan<- []byte {
	bufferChannel := make(chan []byte, p.maxBufferedMessages)
	go p.writeMessagesUntilClose(writer, bufferChannel)
	return bufferChannel
}

func (p *protobufMessageWriter) writeMessagesUntilClose(writer io.Writer, bufferChannel <-chan []byte) {
	msgLengthBuffer := make([]byte, 8)

	// When we are in discard mode, all messages received are just discarded
	discardMode := false

	for msg := range bufferChannel {
		if !discardMode {
			// Write the length
			n := binary.PutUvarint(msgLengthBuffer, uint64(len(msg)))
			msgLengthBuf := msgLengthBuffer[:n]

			// Assemble complete buffer to write
			sendBuf := append(msgLengthBuf, msg...)
			if _, err := writer.Write(sendBuf); err != nil {
				discardMode = true
			}
		}
	}
}
