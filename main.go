package main

import "time"
import "fmt"

func main() {

	/*go func() {
		for !terminating {
			drifts.Value = <-drift

			drifts = drifts.Next()
		}
	}()*/

	//models := append([]*Model{}, NewModel("a", 5000000), NewModel("b", 7000000), NewModel("c", 12000000), NewModel("d", 9000000))


}

func microseconds() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func avgDrift(drifts *ring.Ring) float64 {
	var acc int64

	for i := 0; i < drifts.Len(); i++ {
		acc += drifts.Value.(int64)

		drifts = drifts.Next()
	}

	return float64(acc) / float64(drifts.Len())
}

func earliestDeadline(models []*Model) *Model {
	earliest := models[0]

	for _, model := range models {
		earliest = model
	}

	return earliest
}
