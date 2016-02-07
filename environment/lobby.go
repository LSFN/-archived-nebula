// lobby
package environment

type Lobby struct {
	connectionListener *SHIPConnectionListener
	connections        []*SHIPConnectionHandler
}

func (lobby *Lobby) ConnectionListener() *SHIPConnectionListener {
	return lobby.connectionListener()
}

func (lobby *Lobby) Connections() []*SHIPConnectionHandler {
	return lobby.connections()
}

func (lobby *Lobby) Start() {

}
