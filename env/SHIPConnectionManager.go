// SHIPConnectionManager
package env

import (
	"fmt"
	"net"
)

const (
	LISTEN_FAILED = iota
	SHIP_CONNECTED
	SHIP_DISCONNECTED
)

type SHIPConnectionManagementSignal struct {
	sigType  int
	err      error
	messages chan SHIPMessage
}

type SHIPConnectionManager struct {
	control chan SHIPConnectionManagementSignal
}

func (cm *SHIPConnectionManager) Start() {
	cm.control = make(chan SHIPConnectionManagementSignal)
	listener, err := net.Listen("tcp", ":39461")
	if err != nil {
		cm.control <- SHIPConnectionManagementSignal{sigType: LISTEN_FAILED, err: err}
	} else {
		fmt.Println("Listening")
		go cm.awaitConnections(listener)
	}
}

func (cm *SHIPConnectionManager) awaitConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Drop connection
			conn.Close()
			continue
		}
		handler := new(SHIPConnectionHandler)
		go handler.HandleSHIPConnection(conn)
		cm.control <- SHIPConnectionManagementSignal{sigType: SHIP_CONNECTED, messages: handler.IncomingMessages}
	}
}
