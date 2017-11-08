package main

import "sync"
import "time"
import "fmt"
import "container/ring"

func main() {
	terminating := false
	clockRate := int64(1000000)

	drift := make(chan int64)
	drifts := ring.New(10)

	go func() {
		for !terminating {
			drifts.Value = <-drift

			drifts = drifts.Next()
		}
	}()

	var wg sync.WaitGroup

	models := append([]*Model{}, NewModel("a", 5000000), NewModel("b", 7000000), NewModel("c", 12000000), NewModel("d", 9000000))

	for !terminating {
		//fmt.Printf("scheduling @%d...\n", microseconds())

		time.Sleep(time.Duration(clockRate) * time.Microsecond)

		for _, model := range models {
			wg.Add(1)

			now := microseconds()

			nextCycle := model.lastCycle + model.cycleTime
			nextTick := now + clockRate

			if nextCycle < nextTick {

				//avg := avgDrift(drifts)
				avg := 0.1
				fmt.Printf("avg drift: %f\n", avg)

				go func(model *Model, drift chan int64) {
					time.Sleep(time.Duration(nextCycle - microseconds()))

					localDrift := microseconds() - nextCycle

					if model.lastCycle != 0 {
						drift <- localDrift
					}

					model.lastCycle = microseconds()

					fmt.Printf("planned: %d, scheduled: %d, drift: %d\n", nextCycle, microseconds(), localDrift)

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

func avgDrift(drifts *ring.Ring) float64 {
	var acc int

	for i := 0; i < drifts.Len(); i++ {
		acc += drifts.Value.(int)

		drifts = drifts.Next()
	}

	return float64(acc / drifts.Len())
}

func earliestDeadline(models []*Model) *Model {
	earliest := models[0]

	for _, model := range models {
		earliest = model
	}

	return earliest
}
