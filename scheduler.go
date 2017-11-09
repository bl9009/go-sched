import "container/ring"

type Scheduler struct {
    scheduleDrift chan int64

    clockRate int64
    models []*Model

    drifts *ring.Ring
    terminating bool
}

func (scheduler *Scheduler) NewScheduler(models []*Model, clockRate int64) {
    scheduler.clockRate = clockRate
    scheduler.models = models

    scheduler.terminating = false

    scheduler.scheduleDrift = make(chan int64)

    for i := 0; i < scheduler.drifts.Len(); i++ {
		scheduler.drifts.Value = int64(0)

		scheduler.drifts = scheduler.drifts.Next()
	}
}

func (scheduler *Scheduler) ScheduleSync() {
    for !scheduler.terminating {
        //fmt.Printf("scheduling @%d...\n", microseconds())

        time.Sleep(time.Duration(scheduler.clockRate) * time.Microsecond)

        for _, model := range scheduler.models {
            now := microseconds()

            nextCycle := model.lastCycle + model.cycleTime - int64(avgDrift(scheduler.drifts))
            nextTick := now + scheduler.clockRate

            if nextCycle < nextTick {

                go func(model *Model, drift chan int64) {
                    time.Sleep(time.Duration(nextCycle - microseconds()))

                    localDrift := microseconds() - nextCycle

                    if model.lastCycle != 0 {
                        // this is blocking and might add latency
                        // as channel is synchronized.
                        // preferably replace by queued channel
                        drift <- localDrift 
                    }

                    model.lastCycle = microseconds()

                    fmt.Printf("planned: %d, scheduled: %d, drift: %d\n", nextCycle, microseconds(), localDrift)

                    model.Run()
                }(model, scheduler.scheduleDrift)
            }
        }
    }
}

func (scheduler *Scheduler) ScheduleAsync() {
    go scheduler.ScheduleSync()
}
