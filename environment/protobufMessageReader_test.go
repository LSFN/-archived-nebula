package environment

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		maxBytesPerMessage         uint64
		expectedMaxBytesPerMessage uint64
		expectedBytesForLength     uint8
		maxBufferedMessages        uint
		expectedBufferedMessages   uint
	}{
		{0, 255, 1, 10, 10},
		{255, 255, 1, 0, 1},
		{255, 255, 1, 10, 10},
		{256, 256, 2, 10, 10},
		{(1 << 16) - 1, (1 << 16) - 1, 2, 10, 10},
		{(1 << 16), (1 << 16), 3, 10, 10},
		{(1 << 24) - 1, (1 << 24) - 1, 3, 10, 10},
		{(1 << 24), (1 << 24), 4, 10, 10},
		{(1 << 32) - 1, (1 << 32) - 1, 4, 10, 10},
		{(1 << 32), (1 << 32), 5, 10, 10},
	}

	for _, test := range tests {
		pmr := NewProtobufMessageReader(test.maxBytesPerMessage, test.maxBufferedMessages)
		if pmr.maxBytesPerMessage != test.expectedMaxBytesPerMessage {
			t.Fatalf("maxBytesPerMessage was %d, expected %d", pmr.maxBytesPerMessage, test.expectedMaxBytesPerMessage)
		}
		if pmr.bytesForLength != test.expectedBytesForLength {
			t.Fatalf("bytesForLength was %d, expected %d", pmr.bytesForLength, test.expectedBytesForLength)
		}
		if pmr.maxBufferedMessages != test.expectedBufferedMessages {
			t.Fatalf("Max bytes was %d, expected %d", pmr.maxBufferedMessages, test.maxBufferedMessages)
		}
	}
}

func TestReadMessages(t *testing.T) {
	pmr := NewProtobufMessageReader((1<<16)-1, 1)

	message := "Hello, World!"
	buf := new(bytes.Buffer)
	length := uint16(len(message))
	fmt.Println("Message length is", length)
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		fmt.Println("error", err)
	}
	fmt.Printf("Buffer: % x\n", buf)
	buf.WriteString(message)
	fmt.Printf("Buffer: % x\n", buf)
	bufferChannel := make(chan []byte)

	pmr.readMessagesUntilError(buf, bufferChannel)
	rawMessageBuffer := <-bufferChannel
	convertedMessageBuffer := string(rawMessageBuffer)
	if convertedMessageBuffer != message {
		t.Fatalf("convertedMessageBuffer was %s, expected %s", convertedMessageBuffer, message)
	}

}
