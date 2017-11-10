package main

import ()

func main() {

	/*go func() {
		for !terminating {
			drifts.Value = <-drift

			drifts = drifts.Next()
		}
	}()*/

	models := append([]*Model{}, NewModel("a", 5000000), NewModel("b", 7000000), NewModel("c", 12000000), NewModel("d", 9000000))

	scheduler := NewScheduler(models, 1000000)

	scheduler.ScheduleAsync()
}
