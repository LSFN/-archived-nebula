// env
package environment

import (
	"fmt"
)

// Default settings
const (
	LISTENING_PORT = 39461
)

// Game States
type GameStateName int

const (
	UNINITIALISED = iota
	LOBBY
	SETUP
	PLAY
	CLEANUP
)

type GameState interface {
	ConnectionListener() *SHIPConnectionListener
	Connections() []*SHIPConnectionHandler
	Start()
}

type ENV struct {
}

func (env *ENV) start() {
	var currentGameState *GameState
	currentGameStateName := UNINITIALISED
	nextGameStateName := LOBBY

	// Things to pass to the following state
	var connectionListener *SHIPConnectionListener
	var shipConnections []*SHIPConnectionHandler

	for {
		switch currentGameStateName {
		case UNINITIALISED:
			connectionListener = new(SHIPConnectionListener)
			connectionListener.Start(LISTENING_PORT)
		case LOBBY:
		case SETUP:
		case PLAY:
		case CLEANUP:
		}

		currentGameStateName = nextGameStateName

		switch currentGameStateName {
		case LOBBY:

		case SETUP:
		case PLAY:
		case CLEANUP:
		}
	}
}
