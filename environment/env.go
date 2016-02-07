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
type GamePhase int

const (
	PHASE_UNINITIALISED = iota
	PHASE_LOBBY
	PHASE_SETUP
	PHASE_PLAY
	PHASE_CLEANUP
)

type ENV struct {
}

func (env *ENV) Start() {
	fmt.Println("Starting ENV")
	gameState := new(GameState) // Might just put GameState stuff in ENV struct.
	currentGamePhase := PHASE_UNINITIALISED
	nextGamePhase := PHASE_LOBBY

	var lobby *Lobby

	for {
		switch currentGamePhase {
		case PHASE_UNINITIALISED:
			gameState.connectionManager = new(SHIPConnectionManager)
			gameState.connectionManager.Start(LISTENING_PORT)
		case PHASE_LOBBY:
		case PHASE_SETUP:
		case PHASE_PLAY:
		case PHASE_CLEANUP:
		}

		currentGamePhase = nextGamePhase

		switch currentGamePhase {
		case PHASE_LOBBY:
			lobby.Start(gameState)
		case PHASE_SETUP:
		case PHASE_PLAY:
		case PHASE_CLEANUP:
		}
	}
}
