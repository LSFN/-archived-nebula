// lobby
package environment

import (
	"github.com/LSFN/seprotocol"

	"github.com/pborman/uuid"
)

type Lobby struct {
	gameState *GameState
}

type ShipInfo struct {
	shipServerID string
	shipName     string
}

func (lobby *Lobby) Start(gameState *GameState) {
	// TODO start using infromation from a game coming from the cleanup phase
	lobby.gameState = gameState
	lobby.gameState.shipInfo = make(map[string]*ShipInfo)
	go lobby.handleConnectionEvents(gameState.connectionManager.info)
}

func (lobby *Lobby) handleConnectionEvents(connEventChan chan SCMInfo) {
	for connEvent := range connEventChan {
		switch connEvent.msgType {
		case SCM_LISTEN_FAILED:
			panic("SHIPConnectionManager stopped listening.")
			// TODO something sensible
		case SCM_SHIP_CONNECTING:
			// Start listening to new connections
			go lobby.listenToShipServer(lobby.gameState.connectionManager.connections[connEvent.connectionID])
		case SCM_SHIP_DISCONNECTING:

		}
	}
}

func (lobby *Lobby) listenToShipServer(messenger *ShipServerMessenger) {
	hasJoined := false

	for msg := range messenger.inbound {
		if !hasJoined {
			// The ship server must send a join request first
			if msg.JoinRequest == nil {
				// The ship didn't send a join message, so send a failure message
				messenger.outbound <- &seprotocol.Downstream{
					JoinResponse: &seprotocol.JoinResponse{
						GameVersion: NEBULA_VERSION,
						JoinSuccess: false,
						GamePhase:   seprotocol.JoinResponse_LOBBY,
					},
				}

				// then disconnect it.
				close(messenger.outbound)
				break
			} else {
				// Ships cannot reconnect per-se in the lobby phase so they are always issued a new UUID
				shipInfo := &ShipInfo{shipServerID: uuid.New(), shipName: msg.SetShipName}

				// Send the join response
				messenger.outbound <- &seprotocol.Downstream{
					JoinResponse: &seprotocol.JoinResponse{
						GameVersion:  NEBULA_VERSION,
						JoinSuccess:  true,
						GamePhase:    seprotocol.JoinResponse_LOBBY,
						ShipServerID: shipInfo.shipServerID,
					},
				}

				// Send a lobby join notification to all other ship servers
				lobby.gameState.connectionManager.sendToAll(&seprotocol.Downstream{
					LobbyMembership: &seprotocol.LobbyMembership{
						InfoType: seprotocol.LobbyMembership_JOIN,
						LobbyMembers: []*seprotocol.LobbyMembership_LobbyMemberInfo{
							&seprotocol.LobbyMembership_LobbyMemberInfo{
								ShipServerID: shipInfo.shipServerID,
								ShipName:     shipInfo.shipName,
							},
						},
					},
				})

				// Acknowledge join success in local state
				lobby.gameState.shipInfo[shipInfo.shipServerID] = shipInfo
				hasJoined = true

				// Send a lobby membership message to the newly connecting ship server
				messenger.outbound <- &seprotocol.Downstream{
					LobbyMembership: lobby.makeLobbyMembershipMessage(),
				}
			}
		} else {
			// TODO other lobby messages
		}
	}
}

func (lobby *Lobby) makeLobbyMembershipMessage() *seprotocol.LobbyMembership {
	lobbyMembers := make([]*seprotocol.LobbyMembership_LobbyMemberInfo, len(lobby.gameState.shipInfo))
	i := 0
	for _, shipInfo := range lobby.gameState.shipInfo {
		lobbyMembers[i] = &seprotocol.LobbyMembership_LobbyMemberInfo{
			ShipServerID: shipInfo.shipServerID,
			ShipName:     shipInfo.shipName,
		}
		i++
	}
	return &seprotocol.LobbyMembership{
		InfoType:     seprotocol.LobbyMembership_COMPLETE_LIST,
		LobbyMembers: lobbyMembers,
	}
}
