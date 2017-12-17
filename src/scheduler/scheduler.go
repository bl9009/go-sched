package scheduler

import (
	"errors"
	"time"
)

type Scheduler struct {
	clockRate int64 // [microseconds]
	tasks     []*Task

	dispatcher  *Dispatcher
	terminating bool
	running     bool
}

func NewScheduler(tasks []*Task, clockRate int64) *Scheduler {
	scheduler := new(Scheduler)

	scheduler.clockRate = clockRate
	scheduler.tasks = tasks

	scheduler.dispatcher = NewDispatcher()

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
	for _, task := range scheduler.tasks {

		scheduler.dispatch(task, async)
	}
}

func (scheduler *Scheduler) dispatch(task *Task, async bool) {
	nextTick := scheduler.nextTick()
	nextCycle := scheduler.nextCycle(task)

	if nextCycle < nextTick {
		if async {
			scheduler.dispatcher.DispatchAsync(task, nextCycle)
		} else {
			scheduler.dispatcher.DispatchSync(task, nextCycle)
		}
	}
}

func (scheduler *Scheduler) nextTick() int64 {
	return nowMicroseconds() + scheduler.clockRate
}

func (scheduler *Scheduler) nextCycle(task *Task) int64 {
	return task.lastCycle + task.cycleTime - scheduler.dispatcher.AvgDrift()
}

func (scheduler *Scheduler) waitUntilTick() {
	time.Sleep(time.Duration(scheduler.clockRate) * time.Microsecond)
}

func nowMicroseconds() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}
