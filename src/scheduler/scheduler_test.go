package scheduler

import (
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	testClockRate := int64(10000)
	testModels := append([]*Model{}, NewModel("a", 50000), NewModel("b", 70000), NewModel("c", 120000))

	scheduler := NewScheduler(testModels, testClockRate)

	if scheduler.clockRate != testClockRate {
		t.Errorf("clockRate not set properly: %d != %d", scheduler.clockRate, testClockRate)
	}

	for _, testModel := range testModels {

		found := false

		for _, model := range scheduler.models {
			if testModel == model {
				found = true
			}
		}

		if !found {
			t.Errorf("models not set properly")
		}
	}

	if scheduler.drifts.Len() != 10 {
		t.Errorf("not sufficient elements in drifts buffer: %d", scheduler.drifts.Len())
	}

	scheduler.drifts.Do(func(value interface{}) {
		if scheduler.drifts.Value != int64(0) {
			t.Errorf("drifts not set properly")
		}
	})
}

func TestScheduleAsync(t *testing.T) {

}

func TestScheduleSync(t *testing.T) {

}

func TestSchedule(t *testing.T) {
	scheduler := NewScheduler([]*Model{}, 50000)
	exitChan := make(chan bool)
	startChan := make(chan bool)

	if scheduler.running {
		t.Errorf("scheduler running flag is true although scheduler has not been started yet")
	}

	go scheduler.schedule(false, startChan, exitChan)

	<-startChan

	if !scheduler.running {
		t.Errorf("scheduler running flag is false although scheduler was started")
	}

	scheduler.Terminate()

	<-exitChan

	if scheduler.running {
		t.Errorf("scheduler running flag is true although scheduler has been stopped")
	}
}

func TestTick(t *testing.T) {

}

func TestDispatch(t *testing.T) {

}

func TestNextTick(t *testing.T) {

}

func TestNextModelCycle(t *testing.T) {

}

func TestWaitUntilTick(t *testing.T) {

}

func TestNowMicroseconds(t *testing.T) {
	nowTest := nowMicroseconds()
	now := time.Now().UnixNano() / int64(time.Microsecond)

	if nowTest != now {
		t.Errorf("divergent test timestamp: %d", nowTest)
	}
}

func TestAvgDrift(t *testing.T) {

}
