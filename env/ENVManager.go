// ENVManager
package env

import (
	"fmt"
)

type ENVManager struct {
	
}

func (em *ENVManager) Start(terminate chan int) {
	connectionManager := new(SHIPConnectionManager)
	connectionManager.Start()
	for {
		msg := <-connectionManager.control
		fmt.Println(msg)
	}
}
