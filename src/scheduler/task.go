package scheduler

const (
	INITIALIZED = 0
	SCHEDULED   = 1
	RUNNING     = 2
	SLEEPING    = 3
	FINISHED    = 4
)

type Task struct {
	name      string
	lastCycle int64 // [ms]
	cycleTime int64 // [ms]
	deadline  int   // [ms]
	state     int
}

func NewTask(name string, cycleTime int64) *Task {
	task := new(Task)

	task.name = name
	task.cycleTime = cycleTime

	task.state = INITIALIZED

	return task
}

func (t *Task) Schedule() {
	t.state = SCHEDULED
}

func (t *Task) Run() {
	t.state = RUNNING

	t.lastCycle = nowMicroseconds()

	t.Execute()

	t.state = SLEEPING
}

func (t *Task) Execute() {

}
