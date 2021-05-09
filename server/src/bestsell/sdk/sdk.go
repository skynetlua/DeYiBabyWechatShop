package sdk

import (
	"fmt"
)

func StartServer(ch *chan bool) {
	go func() {
		fmt.Println("StartServer sdk")
			//OnBmap()
		(*ch) <- true
	}()
	<-(*ch)
}
