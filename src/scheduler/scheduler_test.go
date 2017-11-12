package scheduler

import (
    "testing"
    "time"
)

func TestNowMicroseconds(t *testing.T) {
    //scheduler:= NewScheduler([]*Model{}, 0)
    nowTest := nowMicroseconds()
    now := time.Now().UnixNano() / int64(time.Microsecond)

    if nowTest != now {
        t.Errorf("divergent test timestamp: ")
    }
}
