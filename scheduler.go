package main

import (
    "time"
    "container/ring"
)

type Scheduler struct {
    clockRate int64
    models []*Model

    cycleDrift chan int64

    dispatcher *Dispatcher
    drifts *ring.Ring
    terminating bool
}

func NewScheduler(models []*Model, clockRate int64) *Scheduler {
    scheduler := new(Scheduler)

    scheduler.clockRate = clockRate
    scheduler.models = models

    scheduler.cycleDrift = make(chan int64)

    scheduler.dispatcher = NewDispatcher(scheduler.cycleDrift)

    for i := 0; i < scheduler.drifts.Len(); i++ {
		scheduler.drifts.Value = int64(0)

		scheduler.drifts = scheduler.drifts.Next()
	}

    scheduler.terminating = false

    return scheduler
}

func (scheduler *Scheduler) ScheduleAsync() {
    go scheduler.schedule(true)
}

func (scheduler *Scheduler) ScheduleSync() {
    scheduler.schedule(false)
}

func (scheduler *Scheduler) schedule(async bool) {
    for !scheduler.terminating {
        scheduler.waitUntilTick()

        scheduler.tick(async)
    }
}

func (scheduler *Scheduler) tick(async bool) {
    for _, model := range scheduler.models {

        scheduler.dispatch(model, async)
    }
}

func (scheduler *Scheduler) dispatch(model *Model, async bool) {
    nextTick := scheduler.nextTick()
    nextModelCycle := scheduler.nextModelCycle(model)

    if nextModelCycle < nextTick {
        if async {
            scheduler.dispatcher.DispatchAsync(model, nextModelCycle)
        } else {
            scheduler.dispatcher.DispatchSync(model, nextModelCycle)
        }
    }
}

func (scheduler *Scheduler) nextTick() int64 {
    return nowMicroseconds() + scheduler.clockRate
}

func (scheduler *Scheduler) nextModelCycle(model *Model) int64 {
    return model.lastCycle + model.cycleTime - int64(avgDrift(scheduler.drifts))
}

func (scheduler *Scheduler) waitUntilTick() {
    time.Sleep(time.Duration(scheduler.clockRate) * time.Microsecond)
}

func nowMicroseconds() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func avgDrift(drifts *ring.Ring) int64 {
	var acc int64

	for i := 0; i < drifts.Len(); i++ {
		acc += drifts.Value.(int64)

		drifts = drifts.Next()
	}

	return acc / int64(drifts.Len())
}
