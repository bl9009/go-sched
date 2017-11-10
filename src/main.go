package main

import (
    "scheduler"
)

func main() {

	/*go func() {
		for !terminating {
			drifts.Value = <-drift

			drifts = drifts.Next()
		}
	}()*/

	models := append([]*scheduler.Model{}, scheduler.NewModel("a", 5000000), scheduler.NewModel("b", 7000000), scheduler.NewModel("c", 12000000), scheduler.NewModel("d", 9000000))

	scheduler := scheduler.NewScheduler(models, 1000000)

	scheduler.ScheduleAsync()
}
