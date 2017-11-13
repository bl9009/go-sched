package scheduler

import (
	"fmt"
	"time"
)

const (
	INITIALIZED = 0
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
	state     int
}

func NewModel(name string, cycleTime int64) *Model {
	model := new(Model)

	model.name = name
	model.cycleTime = cycleTime

	model.state = INITIALIZED

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
