// DownstreamConnectionManager
package environment

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/LSFN/seprotocol"

	"github.com/blang/semver"
)

const (
	VERSION_HANDSHAKE_TIMEOUT = 1 // in seconds
	NEBULA_PROTOCOL_VERSION   = "0.1.0"
	ACCETPED_VERSION_RANGE    = ">=0.1.0 <0.2.0"
)

type SCMInfoType int

const (
	SCM_LISTEN_FAILED = iota
	SCM_SHIP_CONNECTING
	SCM_SHIP_DISCONNECTING
)

type SCMInfo struct {
	msgType      SCMInfoType
	err          error
	connectionID string
}

type ShipServerMessenger struct {
	inbound  <-chan *seprotocol.Upstream
	outbound chan<- *seprotocol.Downstream
}

type DownstreamConnectionManager struct {
	info        chan SCMInfo
	connections map[string]*ShipServerMessenger
}

func (cm *DownstreamConnectionManager) Start(port uint16) {
	cm.info = make(chan SCMInfo)

	go cm.listen(port)
}

func (cm *DownstreamConnectionManager) listen(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		cm.info <- SCMInfo{msgType: SCM_LISTEN_FAILED, err: err}
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			cm.info <- SCMInfo{msgType: SCM_LISTEN_FAILED, err: err}
			return
		}
		handler := new(DownstreamConnectionHandler)
		handler.Start(conn)
		go func() {
			err := cm.performVersionHandshake(handler)
			if err != nil {
				fmt.Println(err)
			} else {
				messenger := cm.manageConnection(handler)
				cm.connections[handler.id] = messenger
				cm.info <- SCMInfo{msgType: SCM_SHIP_CONNECTING, connectionID: handler.id}
			}
		}()
	}
}

func (cm *DownstreamConnectionManager) performVersionHandshake(handler *DownstreamConnectionHandler) error {
	// Perform a version handshake with the ship server.
	// The ship server first sends its version
	var handshakeError error
	select {
	case <-time.After(time.Second * 1):
		// The ship server hasn't immediately provided its protocol version
		// Close the connection
		close(handler.outboundMessages)
		handshakeError = errors.New("Ship server didn't perform protocol version handshake.")
	case msg := <-handler.inboundMessages:
		// Read the version field
		version, err := semver.Parse(msg.ProtocolVersion)
		if err != nil {
			close(handler.outboundMessages)
			handshakeError = errors.New("Ship server didn't provide a valid protocol version.")
		}
		versionRange, err := semver.ParseRange(ACCETPED_VERSION_RANGE)
		if err != nil {
			panic("Constant ACCETPED_VERSION_RANGE is not a valid SemVer range! Blame the dev.")
		}
		if !versionRange(version) {
			close(handler.outboundMessages)
			handshakeError = errors.New("Ship server version (" + version.String() + ") doesn't satisfy required version range (" + ACCETPED_VERSION_RANGE + ").")
		}
	}

	// Regardless of failure, send our protocol version
	handler.outboundMessages <- &seprotocol.Downstream{ProtocolVersion: NEBULA_PROTOCOL_VERSION}

	// If there was an error, close the connection
	if handshakeError != nil {
		close(handler.outboundMessages)
	}

	// Returns nil if a successful handshake was performed, error otherwise
	return handshakeError
}

func (cm *DownstreamConnectionManager) manageConnection(handler *DownstreamConnectionHandler) *ShipServerMessenger {
	messenger := new(ShipServerMessenger)
	inbound := make(chan *seprotocol.Upstream)
	messenger.inbound = inbound
	outbound := make(chan *seprotocol.Downstream)
	messenger.outbound = outbound

	// Inbound messages
	go func() {
		for msg := range handler.inboundMessages {
			inbound <- msg
		}
		close(inbound)
		delete(cm.connections, handler.id)
		cm.info <- SCMInfo{msgType: SCM_SHIP_DISCONNECTING, connectionID: handler.id}
	}()

	// Outbound messages
	go func() {
		for msg := range outbound {
			handler.outboundMessages <- msg
		}
		close(handler.outboundMessages)
	}()

	return messenger
}

func (cm *DownstreamConnectionManager) sendToAll(message *seprotocol.Downstream) {
	for _, connection := range cm.connections {
		connection.outbound <- message
	}
}
