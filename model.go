package main

import "fmt"
import "time"

const (
	INITILAZIED = 0
	SCHEDULED   = 1
	RUNNING     = 2
	SLEEPING    = 3
	FINISHED    = 4
)

type Model struct {
	name      string
	lastCycle int64 // [ms]
	cycleTime int64 // [ms]
	deadline  int   // [ms]
	state     int   //:= INITIALIZED // 0: initialized, 1: scheduled, 2: running, 3: finished
}

func NewModel(name string, cycleTime int64) *Model {
	model := new(Model)

	model.name = name
	model.cycleTime = cycleTime

	return model
}

func (m Model) Schedule() {
	m.state = SCHEDULED
}

func (m Model) Run() {
	m.state = RUNNING

    m.lastCycle = nowMicroseconds()

	fmt.Printf("model %s @%p\n", m.name, &m)

	time.Sleep(10000 * time.Millisecond)

	m.state = SLEEPING
}
