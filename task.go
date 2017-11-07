package main

import "fmt"

const (
	INITILAZIED = 0
	SCHEDULED   = 1
	RUNNING     = 2
	FINISHED    = 3
)

type Task struct {
	name   string
	status int //:= INITIALIZED // 0: initialized, 1: scheduled, 2: running, 3: finished
}

func (t Task) Schedule() {
	t.status = SCHEDULED
}

func (t Task) Do() {
	t.status = RUNNING

	fmt.Printf("task %s @%p\n", t.name, &t)

	t.status = FINISHED
}
