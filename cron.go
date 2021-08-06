// Copyright ©2021 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sun

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var _ cron.Schedule = (*Schedule)(nil)

// Event is the set of solar events.
type Event int

const (
	Sunrise Event = iota - 1
	Noon
	Sunset
)

func (e Event) String() string {
	switch e {
	case Sunrise:
		return "@sunrise"
	case Noon:
		return "@noon"
	case Sunset:
		return "@sunset"
	default:
		panic(fmt.Sprintf("invalid event: %d", e))
	}
}

// Schedule is a cron.Schedule that schedules solar-time events.
type Schedule struct {
	// Event is the event type being scheduled.
	Event Event
	// Offset is the relative time offset relative
	// to the solar-time event.
	Offset time.Duration

	// Lat and Lon are the latitude and longitude
	// for the scheduled event.
	Lat, Lon float64

	// Location overrides the schedule location.
	Location *time.Location
}

// Next returns the next time the receiver's event occurs.
func (s *Schedule) Next(t time.Time) time.Time {
	next := event(t, s.Lat, s.Lon, s.Event).Add(s.Offset)
	if t.Before(next) {
		return next
	}
	return event(t.Add(24*time.Hour), s.Lat, s.Lon, s.Event).Add(s.Offset)
}

// event returns the time of the specified event at the given
// location on the given date.
func event(date time.Time, lat, lon float64, e Event) time.Time {
	date = time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location())
	rise, noon, set := Times(date, lat, lon)
	switch e {
	case Sunrise:
		return rise
	case Noon:
		return noon
	case Sunset:
		return set
	default:
		panic(fmt.Sprintf("invalid event: %d", e))
	}
}

// Parser is a cron spec parser that handles solar event descriptors.
// These are
//  - @sunrise
//  - @noon
//  - @sunset
// Each solar cron spec in in the form
//  @(sunrise|noon|sunset)([+-]duration)? lat lon
type Parser struct {
	// CronParser is the parser used to handle
	// any non-solar cron specs. If it is nil
	// the standard github.com/robfig/cron/v3
	// parser is used.
	CronParser cron.ScheduleParser
}

var standardCronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// Parse returns a schedule representing the given spec.
func (p Parser) Parse(spec string) (cron.Schedule, error) {
	switch {
	case strings.Contains(spec, "@sunrise"):
		return parseSolarDescriptor(spec, Sunrise)
	case strings.Contains(spec, "@noon"):
		return parseSolarDescriptor(spec, Noon)
	case strings.Contains(spec, "@sunset"):
		return parseSolarDescriptor(spec, Sunset)
	}
	if p.CronParser == nil {
		return standardCronParser.Parse(spec)
	}
	return p.CronParser.Parse(spec)
}

// @(sunrise|noon|sunset)[+-]duration lat lon
func parseSolarDescriptor(spec string, e Event) (*Schedule, error) {
	var err error

	s := &Schedule{Event: e, Location: time.Local}

	if strings.HasPrefix(spec, "TZ=") || strings.HasPrefix(spec, "CRON_TZ=") {
		f := strings.SplitN(spec, " ", 2)
		tz := strings.SplitN(f[0], "=", 2)[1]
		s.Location, err = time.LoadLocation(tz)
		if err != nil {
			return nil, fmt.Errorf("provided bad location %s: %w", tz, err)
		}
		spec = f[1]
	}

	spec = spec[len(e.String()):]

	if strings.HasPrefix(spec, "+") || strings.HasPrefix(spec, "-") {
		f := strings.SplitN(spec, " ", 2)
		s.Offset, err = time.ParseDuration(f[0])
		if err != nil {
			return nil, fmt.Errorf("provided bad offset %s: %w", f[0], err)
		}
		spec = f[1]
	}

	f := strings.Fields(spec)
	if len(f) != 2 {
		return nil, fmt.Errorf("provided bad lat/lon %q", spec)
	}
	s.Lat, err = strconv.ParseFloat(f[0], 64)
	if err != nil {
		return nil, fmt.Errorf("provided bad latitude %q: %w", f[0], err)
	}
	s.Lon, err = strconv.ParseFloat(f[1], 64)
	if err != nil {
		return nil, fmt.Errorf("provided bad longitude %q: %w", f[0], err)
	}

	return s, nil
}
