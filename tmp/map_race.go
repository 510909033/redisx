package tmp

var m = make(map[string]int)

func DemoMapRace() {

	ch := make(chan int, 10)
	go func() {
		//m["aa"] = 111
	}()

	for {
		ch <- 1
		go func() {
			defer func() {
				<-ch
			}()
			_ = m["haha"]
			//log.Println(m["hhhh"])
		}()
	}

}
