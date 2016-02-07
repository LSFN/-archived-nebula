// SHIPConnectionManager
package environment

import (
	"fmt"
	"net"

	"github.com/LSFN/shipenvproto"
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

type SHIPMessenger struct {
	inbound  <-chan *shipenvproto.SHIPtoENV
	outbound chan<- *shipenvproto.ENVtoSHIP
}

type SHIPConnectionManager struct {
	info        chan SCMInfo
	connections map[string]*SHIPMessenger
}

func (cm *SHIPConnectionManager) Start(port uint16) {
	cm.info = make(chan SCMInfo)

	go cm.listen(port)
}

func (cm *SHIPConnectionManager) listen(port uint16) {
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
		handler := new(SHIPConnectionHandler)
		handler.Start(conn)
		messenger := cm.manageConnection(handler)
		cm.connections[handler.id] = messenger
		cm.info <- SCMInfo{msgType: SCM_SHIP_CONNECTING, connectionID: handler.id}
	}
}

func (cm *SHIPConnectionManager) manageConnection(handler *SHIPConnectionHandler) *SHIPMessenger {
	messenger := new(SHIPMessenger)
	inbound := make(chan *shipenvproto.SHIPtoENV)
	messenger.inbound = inbound
	outbound := make(chan *shipenvproto.ENVtoSHIP)
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
