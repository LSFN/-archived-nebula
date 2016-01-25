// SHIPConnectionManager
package env

import (
	"fmt"
	"net"
	"reflect"
)

type SHIPConnectionManagerInputMessageType int

const (
	SCM_TERMINATE = iota
	SCM_LISTEN
)

type SHIPConnectionManagerInputMessage struct {
	messageType SHIPConnectionManagerInputMessageType
	port        uint16
}

type SHIPConnectionManagerOutputMessageType int

const (
	SCM_LISTEN_SUCCESS = iota
	SCM_LISTEN_FAILED
	SCM_SHIP_CONNECTED
	SCM_SHIP_DISCONNECTED
)

type SHIPConnectionManagerOutputMessage struct {
	messageType SHIPConnectionManagerOutputMessageType
	err         error
}

type SHIPConnectionManager struct {
	inputMessageChannel  chan SHIPConnectionManagerInputMessage
	outputMessageChannel chan SHIPConnectionManagerOutputMessage
	chToInt              map[*SHIPConnectionHandler]int
	intToCH              map[int]*SHIPConnectionHandler
	nextCHID             int
	outputToCH           map[chan SHIPConnectionHandlerOutputMessage]*SHIPConnectionHandler
	selectCases          []reflect.SelectCase
}

func (cm *SHIPConnectionManager) Init() {
	cm.inputMessageChannel = make(chan SHIPConnectionManagerInputMessage)
	cm.outputMessageChannel = make(chan SHIPConnectionManagerOutputMessage)
	cm.chToInt = make(map[*SHIPConnectionHandler]int)
	cm.intToCH = make(map[int]*SHIPConnectionHandler)
	cm.outputToCH = make(map[chan SHIPConnectionHandlerOutputMessage]*SHIPConnectionHandler)
	cm.recreateSelectCases()
}

func (cm *SHIPConnectionManager) Start() {
mainLoop:
	for {
		fmt.Println("Waiting for messages on channels", cm.selectCases)
		chosenChannelIndex, valueReceived, recvOk := reflect.Select(cm.selectCases)
		if !recvOk {
			panic("We shouldn't get here because we don't close channels")
		}
		switch cm.selectCases[chosenChannelIndex].Chan.Interface().(type) {
		case chan SHIPConnectionManagerInputMessage:
			ctrlMsg := valueReceived.Interface().(SHIPConnectionManagerInputMessage)
			switch {
			case ctrlMsg.messageType == SCM_TERMINATE:
				// TODO Also close listener
				break mainLoop
			case ctrlMsg.messageType == SCM_LISTEN:
				fmt.Printf("Recieved LISTEN message (port %d)\n", ctrlMsg.port)
				go cm.listen(ctrlMsg.port)
			}
		case chan SHIPConnectionHandlerOutputMessage:

		}
	}
}

func (cm *SHIPConnectionManager) listen(port uint16) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		cm.outputMessageChannel <- SHIPConnectionManagerOutputMessage{messageType: SCM_LISTEN_FAILED, err: err}
		return
	}
	cm.outputMessageChannel <- SHIPConnectionManagerOutputMessage{messageType: SCM_LISTEN_SUCCESS}
	for {
		conn, err := listener.Accept()
		if err != nil {
			cm.outputMessageChannel <- SHIPConnectionManagerOutputMessage{messageType: SCM_LISTEN_FAILED, err: err}
			return
		}
		fmt.Println("New connection")
		handler := new(SHIPConnectionHandler)
		handler.Start(conn)
		chID := cm.nextCHID
		cm.nextCHID++
		cm.chToInt[handler] = chID
		cm.intToCH[chID] = handler
		cm.outputToCH[handler.outputMessageChannel] = handler
		cm.outputMessageChannel <- SHIPConnectionManagerOutputMessage{messageType: SCM_SHIP_CONNECTED}
	}
}

func (cm *SHIPConnectionManager) recreateSelectCases() {
	i := 0
	cm.selectCases = make([]reflect.SelectCase, 1+len(cm.outputToCH))
	cm.selectCases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cm.inputMessageChannel)}
	i++
	for channel := range cm.outputToCH {
		cm.selectCases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)}
		i++
	}
}
