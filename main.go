package main

import "sync"

func main() {
	var wg sync.WaitGroup

	tasks := append([]Task{}, Task{name: "a"}, Task{name: "b"}, Task{name: "c"}, Task{name: "d"})

	for _, task := range tasks {
		wg.Add(1)

		task.Schedule()

		go func(task Task) {
			defer wg.Done()

			task.Do()
		}(task)
	}

	wg.Wait()
}
