package scheduler

import (
	"container/ring"
	"errors"
	"time"
)

type Scheduler struct {
	clockRate int64 // [microseconds]
	models    []*Model

	cycleDrift chan int64

	dispatcher  *Dispatcher
	drifts      *ring.Ring
	terminating bool
	running     bool
}

func NewScheduler(models []*Model, clockRate int64) *Scheduler {
	scheduler := new(Scheduler)

	scheduler.clockRate = clockRate
	scheduler.models = models

	scheduler.cycleDrift = make(chan int64)

	scheduler.dispatcher = NewDispatcher(scheduler.cycleDrift)

	scheduler.drifts = ring.New(10)

	for i := 0; i < scheduler.drifts.Len(); i++ {
		scheduler.drifts.Value = int64(0)

		scheduler.drifts = scheduler.drifts.Next()
	}

	scheduler.terminating = false
	scheduler.running = false

	return scheduler
}

func (scheduler *Scheduler) ScheduleAsync() {
	go scheduler.schedule(true, make(chan bool, 1), make(chan bool, 1))
}

func (scheduler *Scheduler) ScheduleSync() {
	scheduler.schedule(false, make(chan bool, 1), make(chan bool, 1))
}

func (scheduler *Scheduler) Terminate() {
	scheduler.terminating = true
}

func (scheduler *Scheduler) schedule(async bool, start chan bool, exit chan bool) error {
	if !scheduler.running {
		scheduler.running = true

		start <- true
	} else {
		exit <- true

		return errors.New("Scheduler is already running!")
	}

	for !scheduler.terminating {
		scheduler.waitUntilTick()

		scheduler.tick(async)
	}

	scheduler.running = false

	exit <- true

	return nil
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
