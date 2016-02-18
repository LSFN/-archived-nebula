// DownstreamConnectionHandler_test.go
package environment

import (
	"encoding/binary"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/LSFN/seprotocol"
	"github.com/golang/protobuf/proto"
)

type fakeAddr struct {
	network       string
	addressString string
}

func (f fakeAddr) Network() string {
	return f.network
}

func (f fakeAddr) String() string {
	return f.addressString
}

type MockConn struct {
	remoteAddress fakeAddr
	localAddress  fakeAddr

	ServerReader *io.PipeReader
	ServerWriter *io.PipeWriter

	ClientReader *io.PipeReader
	ClientWriter *io.PipeWriter
}

func (c MockConn) Close() error {
	if err := c.ServerWriter.Close(); err != nil {
		return err
	}
	if err := c.ServerReader.Close(); err != nil {
		return err
	}
	return nil
}

func (c MockConn) Read(data []byte) (n int, err error)  { return c.ServerReader.Read(data) }
func (c MockConn) Write(data []byte) (n int, err error) { return c.ServerWriter.Write(data) }

func NewMockConn() MockConn {
	serverRead, clientWrite := io.Pipe()
	clientRead, serverWrite := io.Pipe()

	return MockConn{
		ServerReader: serverRead,
		ServerWriter: serverWrite,
		ClientReader: clientRead,
		ClientWriter: clientWrite,
	}
}

func (conn MockConn) LocalAddr() net.Addr {
	return conn.localAddress
}

func (conn MockConn) RemoteAddr() net.Addr {
	return conn.remoteAddress
}

func (conn MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (conn MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (conn MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestStart(t *testing.T) {
	conn := NewMockConn()

	handler := new(DownstreamConnectionHandler)
	handler.Start(conn)

	if handler.id == "" {
		t.Fatal("Handler's id was nil")
	}
	if handler.inboundMessages == nil {
		t.Fatal("Handler's inboundMessages was nil")
	}
	if handler.outboundMessages == nil {
		t.Fatal("Handler's outboundMessages was nil")
	}
}

func TestInboundMessages(t *testing.T) {
	conn := NewMockConn()

	handler := new(DownstreamConnectionHandler)
	handler.Start(conn)

	message := &seprotocol.Upstream{
		ProtocolVersion: NEBULA_PROTOCOL_VERSION,
	}
	messageBuf, err := proto.Marshal(message)
	if err != nil {
		t.Fatal("Someone wrote this test wrong.")
	}

	lengthBuffer := make([]byte, 0, binary.MaxVarintLen64)
	binary.PutUvarint(lengthBuffer, uint64(len(messageBuf)))
	go func() {
		conn.ClientWriter.Write(lengthBuffer)
		conn.ClientWriter.Write(messageBuf)
	}()

	select {
	case <-time.After(time.Second * 1):
		t.Fatal("Didn't receive message.")
	case msg := <-handler.inboundMessages:
		if !reflect.DeepEqual(message, msg) {
			t.Fatalf("Messages didn't match:\n\texpected: %q\n\tgot: %q", message, msg)
		}
	}
}
