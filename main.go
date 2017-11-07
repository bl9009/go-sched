package main

import "sync"
import "time"
import "fmt"

func main() {
    clockRate := time.Duration(3000)

	var wg sync.WaitGroup

	models := append([]*Model{}, NewModel("a"), NewModel("b"), NewModel("c"), NewModel("d"))

    var next *Model

	for _, model := range models {
		wg.Add(1)

        next = earliestDeadline(models)

		next.Schedule()

		go func(model *Model) {
			defer wg.Done()

			model.Run()
		}(model)

        fmt.Println("scheduling...")

        time.Sleep(clockRate * time.Millisecond)
	}

	wg.Wait()
}

func earliestDeadline(models []*Model) *Model {
    earliest := models[0]

    for _, model := range models {
        earliest = model
    }

    return earliest
}
