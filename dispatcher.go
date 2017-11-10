package main

import (
    "time"
)

type Dispatcher struct {
    cycleDrift chan int64
}

func NewDispatcher(cycleDrift chan int64) *Dispatcher {
    dispatcher := new(Dispatcher)

    dispatcher.cycleDrift = cycleDrift

    return dispatcher
}

func (dispatcher *Dispatcher) DispatchAsync(model *Model, nextModelCycle int64) {
    go dispatcher.DispatchSync(model, nextModelCycle)
}

func (dispatcher *Dispatcher) DispatchSync(model *Model, nextModelCycle int64) {
    waitUntilCycle(nextModelCycle)

    if model.lastCycle != 0 {
        // this is blocking and might add latency
        // as channel is synchronized.
        // preferably replace by queued channel
        dispatcher.cycleDrift <- nowMicroseconds() - nextModelCycle
    }

    model.Run()
}

func waitUntilCycle(nextModelCycle int64) {
    time.Sleep(time.Duration(nextModelCycle - nowMicroseconds()))
}
