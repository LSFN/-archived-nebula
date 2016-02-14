// state
package environment

type GameState struct {
	connectionManager *DownstreamConnectionManager
	shipInfo          map[string]*ShipInfo
}
