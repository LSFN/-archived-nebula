// state
package environment

type GameState struct {
	connectionManager      *DownstreamConnectionManager
	shipInfoByShipServerID map[string]*ShipInfo
	shipInfoByConnectionID map[string]*ShipInfo
}
