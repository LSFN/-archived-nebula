// ENV is LSFN's Environment server.
package main

import (
	"github.com/LSFN/nebula/environment"
)

func main() {
	env := new(environment.Server)
	env.Start()
}
