// Package chronos is a scheduling tool for Go based on:
//  https://github.com/carlescere/scheduler

package chronos

import (
	"time"
)

// Enum of scheduler kind
const (
	Periodic = iota
	Monthly  = iota
	Yearly   = iota
)

const (
	Day  = 24 * time.Hour
	Week =  7 * Day
)

// Auxiliar type that holds the information needed to build the scheduler
type auxiliar struct {
	kind,                        // Enum of scheduler kind
	ammount        int
	notInmediately bool
	start,
	end            time.Time
	unit           time.Duration
}

type Job struct {
	task    func()    // Task to be scheduled
	times,            // Times that the task can be executed, -1 means no limit
	aux     auxiliar  // Holds the values for following API calls
	shedule scheduler // Scheduler to determine when to run the job
}

// Job construction with task assignment
func Schedule(f func()) *Job {
	return &Job{task:f, times:-1}
}

// Defining the number of times
func (j *Job) NTimes(n int) *Job {
	j.times = n
	return j
}

func (j *Job) Once() *Job {
	return j.NTimes(1)
}

func (j *Job) Twice() *Job {
	return j.NTimes(2)
}

// Defining the period size in units
func (j *Job) Every(times ...int) *Job {
	switch len(times) {
	case 0:
		j.aux.ammount := 1
	case 1:
		j.aux.ammount := times[0]
	default:
		panic("Too many arguments in Job.Every()")
	}
	return j
}

// Defining the period's unit duration
func (j *Job) duration(d time.Duration) *Job {
	j.aux.kind = periodicKind
	j.aux.unit = d
	return j
}

func (j *Job) Nanosecond() *Job {
	return j.duration(time.Nanosecond)
}

func (j *Job) Nanoseconds() *Job {
	return j.Nanosecond()
}

func (j *Job) Microsecond() *Job {
	return j.duration(time.Microsecond)
}

func (j *Job) Microseconds() *Job {
	return j.Microsecond()
}

func (j *Job) Millisecond() *Job {
	return j.duration(time.Millisecond)
}

func (j *Job) Milliseconds() *Job {
	return j.Millisecond()
}

func (j *Job) Second() *Job {
	return j.duration(time.Second)
}

func (j *Job) Seconds() *Job {
	return j.Second()
}

func (j *Job) Minute() *Job {
	return j.duration(time.Minute)
}

func (j *Job) Minutes() *Job {
	return j.Minute()
}

func (j *Job) Hour() *Job {
	return j.duration(time.Hour)
}

func (j *Job) Hours() *Job {
	return j.Hour()
}

func (j *Job) Day() *Job {
	return j.duration(Day)
}

func (j *Job) Days() *Job {
	return j.Day()
}

func (j *Job) Week() *Job {
	return j.duration(Week)
}

func (j *Job) Weeks() *Job {
	return j.Week()
}

func (j *Job) Month() *Job {
	j.aux.kind = monthlyKind
	return j
}

func (j *Job) Months() *Job {
	return j.Month()
}

func (j *Job) Year() *Job {
	j.aux.kind = yearlyKind
	return j
}

func (j *Job) Years() *Job {
	return j.Year()
}

// Defining if it should run at the start of the cycle
func(j *Job) NotInmediately() *Job {
	j.aux.notInmediately = true
	return j
}

// Defining the starting and ending times
func(j *Job) At(t time.Time) *Job {
	j.aux.start = t
	return j
}

func(j *Job) In(d time.Duration) *Job {
	return j.At(time.Now().Add(d))
}

func(j *Job) Until(t time.Time) *Job {
	j.aux.end = t
	return j
}

// Scheduling the task
func(j *Job) Done() error {
	var err error
	switch j.aux.kind {
	case periodicKind:
		j.schedule, err = newPeriodic(j.aux.start, j.aux.end, j.aux.ammount,
		                              j.aux.unit, j.aux.notInmediately)
	case monthlyKind:
		j.schedule, err = newMonthly(j.aux.start, j.aux.end, j.aux.ammount,
		                             j.aux.notInmediately)
	case yearlyKind:
		j.schedule, err = newYearly(j.aux.start, j.aux.end, j.aux.ammount,
		                            j.aux.notInmediately)
	}
	if err != nil {
		return err
	}

	// TODO: schedule and execute the task
}