// ENV is LSFN's Environment server.
package main

import (
	"fmt"

	"github.com/LSFN/ENV/env"
)

func main() {
	terminate := make(chan int)
	fmt.Println("Starting")
	envManager := new(env.ENVManager)
	go envManager.Start(terminate)
	<-terminate
	fmt.Println("Exiting")
}
