package module

func StartServer(ch *chan bool) {
	go func() {
		// loadTeamMembers()

		(*ch) <- true
	}()
	<-(*ch)
}