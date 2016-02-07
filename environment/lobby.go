// lobby
package environment

type Lobby struct {
}

func (lobby *Lobby) Start(gameState *GameState) {
	go lobby.handleConnectionEvents(gameState.connectionManager.info)
}

func (lobby *Lobby) handleConnectionEvents(connEventChan chan SCMInfo) {
	for connEvent := range connEventChan {
		switch connEvent.msgType {
		case SCM_LISTEN_FAILED:
			panic("SHIPConnectionManager stopped listening.")
		case SCM_SHIP_CONNECTING:

		case SCM_SHIP_DISCONNECTING:

		}
	}
}
