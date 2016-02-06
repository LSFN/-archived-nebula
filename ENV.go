// ENV is LSFN's Environment server.
package main

import (
	"github.com/LSFN/ENV/environment"
)

// Default settings
const (
	LISTENING_PORT = 39461
)

// Game States
type GameState int

const (
	LOBBY = iota
	SETUP
	PLAY
	CLEANUP
)

type ENV struct {
	gameState              GameState
	shipConnectionListener *environment.SHIPConnectionListener
}

func main() {
	env := new(ENV)
	env.gameState = LOBBY
	env.shipConnectionListener = new(environment.SHIPConnectionListener)
	env.shipConnectionListener.Start(LISTENING_PORT)
}
