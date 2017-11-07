package main

import "sync"
import "time"
import "fmt"

func main() {
	terminating := false
	clockRate := int64(1000000)

    drift := make(chan int64)
    var accDrift int64
    var drifts int

	var wg sync.WaitGroup

	models := append([]*Model{}, NewModel("a", 5000000), NewModel("b", 7000000), NewModel("c", 12000000), NewModel("d", 9000000))

	//var next *Model

    go func() {
        for {
            fmt.Println("waiting for drift...")
            lastDrift := <- drift
            fmt.Printf("received drift %d", lastDrift)
            accDrift += <- drift
            drifts += 1
            fmt.Printf("avg drift: %f\n", float32(accDrift) / float32(drifts))
        }
    }()

	for !terminating {
		//fmt.Printf("scheduling @%d...\n", microseconds())

		time.Sleep(time.Duration(clockRate) * time.Microsecond)

		for _, model := range models {
			wg.Add(1)

			now := microseconds()

            nextCycle := model.lastCycle + model.cycleTime
            nextTick := now + clockRate

			if  nextCycle < nextTick {

				go func(model *Model, drift chan int64) {
                    time.Sleep(time.Duration(nextCycle - microseconds()))

                    localDrift := nextCycle - microseconds()

                    if model.lastCycle != 0 {
                        drift <- localDrift
                    }

                    model.lastCycle = microseconds()

                    //fmt.Printf("planned: %d, scheduled: %d, drift: %d\n", nextCycle, microseconds(), drift)

					defer wg.Done()

					model.Run()
				}(model, drift)
			}
		}
	}

	wg.Wait()
}

func microseconds() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func earliestDeadline(models []*Model) *Model {
	earliest := models[0]

	for _, model := range models {
		earliest = model
	}

	return earliest
}
