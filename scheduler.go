// Package chronos is a scheduling tool for Go based on:
//  https://github.com/carlescere/scheduler

package chronos

import (
	"errors"
	"time"
)

// Enum of scheduler kind
const (
	periodicKind = iota
	monthlyKind  = iota
	yearlyKind   = iota
)

const (
	Day  = 24 * time.Hour
	Week =  7 * Day
)

type scheduler interface {
	// Returns wether there is another event scheduled and the remaining time
	next() (bool, time.Duration)
}

// Auxiliar type that holds the information needed to build the scheduler
type auxiliar struct {
	kind,                        // Enum of scheduler kind
	ammount        int
	notInmediately bool
	start,
	end            time.Time
	unit           time.Duration
}

// Accepts periods in every time unit from ns to weeks, months and years need to
// be considered separately as their length is not constant
type periodic struct {
	start,                 // Start time
	end      time.Time     // End time, zero value means no end
	started  bool          // Internal flag to handle first executions
	ammount  time.Duration // Period
	n        int           // Number of already executed events
}

// Constructor
func newPeriodic(start, end time.Time, ammount int, unit time.Duration, notInmediately bool) (*periodic, error) {
	// Check the input is valid
	if ammount == 0 || unit == 0 {
		return nil, errors.New("0 is not a valid period")
	}
	// If no start time was assigned, use current time
	if start.IsZero() {
		start = time.Now()
	}
	// If notInmediately was called, the starting date should not be returned
	// by periodic.next() call, so we add 1 to the event count to avoid it
	var n int
	if notInmediately {
		n = 1
	}

	return &periodic{start:start, end:end, started:notInmediately,
	                 ammount:time.Duration(ammount*int(unit)), n:n},
	       nil
}

// Auxiliar function that returns the execution time candidate
func (s *periodic) getCandidate() time.Time {
	return s.start.Add(time.Duration(s.n*int(s.ammount)))
}

// Implements scheduler.next()
func (s *periodic) next() (bool, time.Duration) {
	// Calculate the next iteration
	next := s.getCandidate()
	for next.Before(time.Now()) {
		if !s.started {
			break
		}
		s.n++
		next = s.getCandidate()
	}
	if !s.started {
		s.started = true
	}

	// Check if the end date has arrived
	return s.end.IsZero() || next.Before(s.end), next.Sub(time.Now())
}

// Monthly periods need to be considered separately as their length is not
// constant (28-31 days)
type monthly struct {
	start,             // Start time
	end      time.Time // End time, zero value means no end
	started  bool      // Internal flag to handle first executions
	ammount,           // Ammount of months that made up a period
	n        int       // Number of already executed events
}

// Constructor
func newMonthly(start, end time.Time, ammount int, notInmediately bool) (*monthly, error) {
	// Check the input is valid
	if ammount == 0 {
		return nil, errors.New("0 months is not a valid period")
	}
	// If no start time was assigned, use current time
	if start.IsZero() {
		start = time.Now()
	}
	// If notInmediately was called, the starting date should not be returned
	// by periodic.next() call, so we add 1 to the event count to avoid it
	var n int
	if notInmediately {
		n = 1
	}

	return &monthly{start:start, end:end, started:notInmediately,
	                ammount:ammount, n:n},
	       nil
}

func (s *monthly) getCandidate() time.Time {
	res := s.start.AddDate(0, s.n*s.ammount, 0)
	if res.Day() != s.start.Day() {
		res = res.AddDate(0, 0, -res.Day())
	}
	return res
}

// Implements scheduler.next()
func (s *monthly) next() (bool, time.Duration) {
	// Calculate the next iteration
	next := s.getCandidate()
	for next.Before(time.Now()) {
		if !s.started {
			break
		}
		s.n++
		next = s.getCandidate()
	}
	if !s.started {
		s.started = true
	}

	// Check if the end date has arrived
	return s.end.IsZero() || next.Before(s.end), next.Sub(time.Now())
}

// Yearly periods need to be considered separately as
// their length is not constant (365-366 days)
type yearly struct {
	start,             // Start time
	end      time.Time // End time, zero value means no end
	started  bool      // Internal flag to handle first executions
	ammount,           // Ammount of years that made up a period
	n        int       // Number of already executed events
}

// Constructor
func newYearly(start, end time.Time, ammount int, notInmediately bool) (*yearly, error) {
	// Check the input is valid
	if ammount == 0 {
		return nil, errors.New("0 years is not a valid period")
	}
	// If no start time was assigned, use current time
	if start.IsZero() {
		start = time.Now()
	}
	// If notInmediately was called, the starting date should not be returned
	// by periodic.next() call, so we add 1 to the event count to avoid it
	var n int
	if notInmediately {
		n = 1
	}

	return &yearly{start:start, end:end, started:notInmediately,
	               ammount:ammount, n:n},
	       nil
}

func (s *yearly) getCandidate() time.Time {
	res := s.start.AddDate(s.n*s.ammount, 0, 0)
	if res.Day() != s.start.Day() {
		res = res.AddDate(0, 0, -res.Day())
	}
	return res
}

// implements scheduler.next()
func (s *yearly) next() (bool, time.Duration) {
	// Calculate the next iteration
	next := s.getCandidate()
	for next.Before(time.Now()) {
		if !s.started {
			break
		}
		s.n++
		next = s.getCandidate()
	}
	if !s.started {
		s.started = true
	}

	// Check if the end date has arrived
	return s.end.IsZero() || next.Before(s.end), next.Sub(time.Now())
}
