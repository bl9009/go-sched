package main

import "sync"
import "time"
import "fmt"

func main() {
	terminating := false
	clockRate := int64(1000000)

	var wg sync.WaitGroup

	models := append([]*Model{}, NewModel("a", 5000000), NewModel("b", 7000000), NewModel("c", 12000000), NewModel("d", 9000000))

	//var next *Model

	for !terminating {
		fmt.Printf("scheduling @%d...\n", microseconds())

		time.Sleep(time.Duration(clockRate) * time.Microsecond)

		for _, model := range models {
			wg.Add(1)

			now := microseconds()

			if model.lastCycle+model.cycleTime < now+clockRate {

				model.lastCycle = now

				go func(model *Model) {
					defer wg.Done()

					model.Run()
				}(model)
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
