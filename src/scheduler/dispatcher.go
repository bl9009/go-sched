package scheduler

import (
	"container/ring"
	"sync"
	"time"
)

type Dispatcher struct {
	mutex  *sync.RWMutex
	drifts *ring.Ring
}

func NewDispatcher() *Dispatcher {
	dispatcher := new(Dispatcher)

	dispatcher.mutex = new(sync.RWMutex)

	dispatcher.drifts = ring.New(10)

	for i := 0; i < dispatcher.drifts.Len(); i++ {
		dispatcher.drifts.Value = int64(0)

		dispatcher.drifts = dispatcher.drifts.Next()
	}

	return dispatcher
}

func (dispatcher *Dispatcher) DispatchAsync(task *Task, nextModelCycle int64) {
	go dispatcher.dispatch(task, nextModelCycle)
}

func (dispatcher *Dispatcher) DispatchSync(task *Task, nextCycle int64) {
	dispatcher.dispatch(task, nextCycle)
}

func (dispatcher *Dispatcher) AvgDrift() int64 {
	dispatcher.mutex.RLock()
	defer dispatcher.mutex.RUnlock()

	var acc int64

	for i := 0; i < dispatcher.drifts.Len(); i++ {
		acc += dispatcher.drifts.Value.(int64)

		dispatcher.drifts = dispatcher.drifts.Next()
	}

	return acc / int64(dispatcher.drifts.Len())
}

func (dispatcher *Dispatcher) dispatch(task *Task, nextCycle int64) {
	waitUntilCycle(nextCycle)

	drift := nowMicroseconds() - nextCycle

	task.Run()

	if task.lastCycle != 0 {
		dispatcher.addDrift(drift)
	}
}

func waitUntilCycle(nextCycle int64) {
	time.Sleep(time.Duration(nextCycle - nowMicroseconds()))
}

func (dispatcher *Dispatcher) addDrift(drift int64) {
	dispatcher.mutex.Lock()
	defer dispatcher.mutex.Unlock()

	dispatcher.drifts.Value = drift

	dispatcher.drifts = dispatcher.drifts.Next()
}
