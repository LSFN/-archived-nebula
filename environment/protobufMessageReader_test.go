package environment

import (
	"bufio"
	"encoding/binary"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		maxBufferedMessages      uint
		expectedBufferedMessages uint
	}{
		{10, 10},
		{0, 1},
	}

	for _, test := range tests {
		pmr := NewProtobufMessageReader(test.maxBufferedMessages)
		if pmr.maxBufferedMessages != test.expectedBufferedMessages {
			t.Fatalf("maxBufferedMessages was %d, expected %d\n", pmr.maxBufferedMessages, test.maxBufferedMessages)
		}
	}
}

func TestReadMessages(t *testing.T) {
	pmr := NewProtobufMessageReader(1)
	pipeReader, pipeWriter := io.Pipe()
	bufferChannel := make(chan []byte)
	go pmr.readMessagesUntilError(bufio.NewReader(pipeReader), bufferChannel)

	message := "Hello, World!"
	rawMessageLength := make([]byte, 8)
	lengthBytes := binary.PutUvarint(rawMessageLength, uint64(len(message)))
	rawMessageLength = rawMessageLength[:lengthBytes]
	if n, err := pipeWriter.Write(rawMessageLength); err != nil {
		t.Fatalf("Couldn't write all of message length to pipeWriter, wrote %d bytes, error %s", n, err)
	}
	if n, err := pipeWriter.Write([]byte(message)); err != nil {
		t.Fatalf("Couldn't write all of message to pipeWriter, wrote %d bytes, error %s", n, err)
	}
	pipeWriter.Close()

	rawMessageBuffer := <-bufferChannel
	convertedMessageBuffer := string(rawMessageBuffer)
	if convertedMessageBuffer != message {
		t.Fatalf("convertedMessageBuffer was %s, expected %s", convertedMessageBuffer, message)
	}

	_, more := <-bufferChannel
	if more {
		t.Fatal("expected bufferChannel to be closed but it isn't")
	}
}
