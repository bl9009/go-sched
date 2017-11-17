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

func (dispatcher *Dispatcher) DispatchAsync(model *Model, nextModelCycle int64) {
	go dispatcher.dispatch(model, nextModelCycle)
}

func (dispatcher *Dispatcher) DispatchSync(model *Model, nextModelCycle int64) {
	dispatcher.dispatch(model, nextModelCycle)
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

func (dispatcher *Dispatcher) dispatch(model *Model, nextModelCycle int64) {
	waitUntilCycle(nextModelCycle)

	drift := nowMicroseconds() - nextModelCycle

	model.Run()

	if model.lastCycle != 0 {
		dispatcher.addDrift(drift)
	}
}

func waitUntilCycle(nextModelCycle int64) {
	time.Sleep(time.Duration(nextModelCycle - nowMicroseconds()))
}

func (dispatcher *Dispatcher) addDrift(drift int64) {
	dispatcher.mutex.Lock()
	defer dispatcher.mutex.Unlock()

	dispatcher.drifts.Value = drift

	dispatcher.drifts = dispatcher.drifts.Next()
}
