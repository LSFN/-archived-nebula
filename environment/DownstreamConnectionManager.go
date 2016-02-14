// DownstreamConnectionManager
package environment

import (
	"fmt"
	"net"

	"github.com/LSFN/seprotocol"
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
		messenger := cm.manageConnection(handler)
		cm.connections[handler.id] = messenger
		cm.info <- SCMInfo{msgType: SCM_SHIP_CONNECTING, connectionID: handler.id}
	}
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
