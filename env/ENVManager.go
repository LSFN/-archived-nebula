// ENVManager
package env

import (
	"fmt"
)

type ControlMessageType int

const (
	TERMINATE = iota
)

type ControlMessage struct {
	messageType ControlMessageType
}

const (
	LISTENING_PORT = 39461
)

type ENVManager struct {
}

func (em *ENVManager) Start(control chan ControlMessage) {
	connectionManager := new(SHIPConnectionManager)
	connectionManager.Init()
	go connectionManager.Start()

	connectionManager.inputMessageChannel <- SHIPConnectionManagerInputMessage{
		messageType: SCM_LISTEN,
		port:        LISTENING_PORT,
	}

mainLoop:
	for {
		select {
		case ctrlMsg := <-control:
			switch ctrlMsg.messageType {
			case TERMINATE:
				break mainLoop
			}
		case scmMsg := <-connectionManager.outputMessageChannel:
			switch scmMsg.messageType {
			case SCM_LISTEN_SUCCESS:
				fmt.Println("Listening on port", LISTENING_PORT)
			case SCM_LISTEN_FAILED:
				fmt.Println("Listen failed:", scmMsg.err)
			case SCM_SHIP_CONNECTED:
				fmt.Println("A new SHIP connected.")
			case SCM_SHIP_DISCONNECTED:
				fmt.Println("A SHIP disconnected.")
			}
		}
	}
}
