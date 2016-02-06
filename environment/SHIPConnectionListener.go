// SHIPConnectionListener
package environment

import (
	"fmt"
	"net"
)

type SCLInfoType int

const (
	SCL_LISTEN_SUCCESS = iota
	SCL_LISTEN_FAILED
	SCL_NEW_SHIP_CONNECTION
)

type SCLInfo struct {
	msgType        SCLInfoType
	err            error
	shipConnection *SHIPConnectionHandler
}

type SHIPConnectionListener struct {
	info chan SCLInfo
}

func (cm *SHIPConnectionListener) Start(port uint16) {
	cm.info = make(chan SCLInfo)

	go cm.listen(port)
}

func (cm *SHIPConnectionListener) listen(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		cm.info <- SCLInfo{msgType: SCL_LISTEN_FAILED, err: err}
		return
	}
	cm.info <- SCLInfo{msgType: SCL_LISTEN_SUCCESS}
	for {
		conn, err := listener.Accept()
		if err != nil {
			cm.info <- SCLInfo{msgType: SCL_LISTEN_FAILED, err: err}
			return
		}
		handler := new(SHIPConnectionHandler)
		handler.Start(conn)
		cm.info <- SCLInfo{msgType: SCL_NEW_SHIP_CONNECTION, shipConnection: handler}
	}
}
