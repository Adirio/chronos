// Package chronos is a scheduling tool for Go based on:
//  https://github.com/carlescere/scheduler

package chronos

import (
	"errors"
	"time"
)

const (
	Day  = 24 * time.Hour
	Week =  7 * Day
)

type ChronosJob struct {
	schedule scheduler
}

type scheduler interface {
	next() (time.Duration, bool, error)
}

// Accepts periods in every time unit from ns to weeks,
// months and years need to be considered separately as
// their length is not constant
type periodic struct {
	start   time.Time
	end     time.Time
	started bool
	ammount int
	unit    time.Duration
}

func (schedule *periodic) next() (time.Duration, bool, error) {
	if schedule.unit == 0 || schedule.ammount == 0 {
		return 0, false, errors.New("0 is not a valid period")
	}
	now := time.Now()
	if !schedule.end.IsZero && now.After(schedule.end) {
		return 0, false, nil
	}
	if !schedule.started {
		schedule.started = true
		if now.After(schedule.start) {
			return 0, true, nil
		}
		return schedule.start.Sub(now), true, nil
	}
	return schedule.ammount * schedule.unit, true, nil
}

// Monthly periods need to be considered separately as
// their length is not constant (28-31 days)
type monthly struct {
	start   time.Time
	end     time.Time
	started bool
	n       int
	ammount int
}

func (schedule *monthly) next() (time.Duration, bool, error) {
	if schedule.ammount == 0 {
		return 0, false, errors.New("0 months is not a valid period")
	}
	now := time.Now()
	if !schedule.end.IsZero && now.After(schedule.end) {
		return 0, false, nil
	}
	if !schedule.started {
		schedule.started = true
		if now.After(schedule.start) {
			return 0, true, nil
		}
		return schedule.start.Sub(now), true, nil
	}
	if now.Before(schedule.start) {
		next := schedule.start
	} else {
		next := schedule.previous.AddDate(0, schedule.n, 0)
		// If the day does not exist for that month time.Time.AddDate
		// normalizes it, so we need to substract the extra days
		if next.Day() != schedule.start.Day() {
			next = next.AddDate(0, 0, -next.Day())
		}
	}
	schedule.n += schedule.ammount
	return next.Sub(now), true, nil
}

// Yearly periods need to be considered separately as
// their length is not constant (365-366 days)
type yearly struct {
	start   time.Time
	end     time.Time
	started bool
	n       int
	ammount int
}

func (schedule *yearly) next() (time.Duration, bool, error) {
	if schedule.ammount == 0 {
		return 0, false, errors.New("0 years is not a valid period")
	}
	now := time.Now()
	if !schedule.end.IsZero && now.After(schedule.end) {
		return 0, false, nil
	}
	if !schedule.started {
		schedule.started = true
		if now.After(schedule.start) {
			return 0, true, nil
		}
		return schedule.start.Sub(now), true, nil
	}
	if now.Before(schedule.start) {
		next := schedule.start
	} else {
		next := schedule.start.AddDate(schedule.n, 0, 0)
		// If the day does not exist for that month time.Time.AddDate
		// normalizes it, so we need to substract the extra days
		if next.Day() != schedule.start.Day() {
			next = next.AddDate(0, 0, -next.Day())
		}
	}
	schedule.n += schedule.ammount
	return next.Sub(now), true, nil
}