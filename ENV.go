// ENV is LSFN's Environment server.
package main

import (
	"github.com/LSFN/ENV/environment"
)

func main() {
	env := new(environment.ENV)
	env.start()
}
