package scheduler

import (
	"math"
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
	// NO TEST REQUIRED: Proxy method only
}

func TestScheduleSync(t *testing.T) {
	// NO TEST REQUIRED: Proxy method only
}

func TestSchedule(t *testing.T) {
	scheduler := NewScheduler([]*Model{}, 50000)
	exitChan := make(chan bool)
	startChan := make(chan bool)

	if scheduler.running {
		t.Errorf("scheduler.running flag is true although scheduler has not been started yet")
	}

	go scheduler.schedule(false, startChan, exitChan)

	<-startChan

	if !scheduler.running {
		t.Errorf("scheduler.running flag is false although scheduler was started")
	}

	err := scheduler.schedule(false, make(chan bool, 1), make(chan bool, 1))

	if err == nil {
		t.Errorf("Scheduler started altough it is already running")
	}

	scheduler.Terminate()

	<-exitChan

	if scheduler.running {
		t.Errorf("scheduler.running flag is true although scheduler has been stopped")
	}
}

func TestTick(t *testing.T) {
	// NO TEST REQUIRED: Structural function without actual logic
}

func TestDispatch(t *testing.T) {
	// NO TEST REQUIRED: Test not required: Is tested in dispatcher tests
}

func TestNextTick(t *testing.T) {
	testClockRate := int64(50000)
	scheduler := NewScheduler([]*Model{}, testClockRate)
	testTolerance := 1000

	testNextTick := nowMicroseconds() + testClockRate

	if math.Abs(float64(testNextTick-scheduler.nextTick())) > float64(testTolerance) {
		t.Errorf("nextTick drifting")
	}
}

func TestNextModelCycle(t *testing.T) {

}

func TestWaitUntilTick(t *testing.T) {
	testClockRate := int64(50000)
	scheduler := NewScheduler([]*Model{}, testClockRate)
	testTolerance := int64(1000)

	testWait := scheduler.clockRate

	start := nowMicroseconds()

	scheduler.waitUntilTick()

	result := nowMicroseconds() - start

	if result < testWait-testTolerance || result > testWait+testTolerance {
		t.Errorf("waitUntilClock drifting: %f", math.Abs(float64(result-testWait)))
	}
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
