// Copyright ©2021 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sun

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var parseTests = []struct {
	spec    string
	want    *Schedule
	wantErr error
}{
	{
		spec:    "@sunrise",
		want:    &Schedule{},
		wantErr: errors.New(`provided bad lat/lon ""`),
	},
	{
		spec:    "@noon",
		want:    &Schedule{},
		wantErr: errors.New(`provided bad lat/lon ""`),
	},
	{
		spec:    "@sunset",
		want:    &Schedule{},
		wantErr: errors.New(`provided bad lat/lon ""`),
	},
	{
		spec: "@sunrise 48.856614 2.3522219",
		want: &Schedule{Event: Sunrise, Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@noon 48.856614 2.3522219",
		want: &Schedule{Event: Noon, Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@sunset 48.856614 2.3522219",
		want: &Schedule{Event: Sunset, Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@sunrise-1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Sunrise, Offset: -(time.Hour + 10*time.Minute), Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@noon-1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Noon, Offset: -(time.Hour + 10*time.Minute), Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@sunset-1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Sunset, Offset: -(time.Hour + 10*time.Minute), Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@sunrise+1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Sunrise, Offset: time.Hour + 10*time.Minute, Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@noon+1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Noon, Offset: time.Hour + 10*time.Minute, Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "@sunset+1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Sunset, Offset: time.Hour + 10*time.Minute, Lat: 48.856614, Lon: 2.3522219, Location: time.Local},
	},
	{
		spec: "TZ=Europe/Paris @sunrise-1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Sunrise, Offset: -(time.Hour + 10*time.Minute), Lat: 48.856614, Lon: 2.3522219, Location: mustLoc("Europe/Paris")},
	},
	{
		spec: "TZ=Europe/Paris @noon-1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Noon, Offset: -(time.Hour + 10*time.Minute), Lat: 48.856614, Lon: 2.3522219, Location: mustLoc("Europe/Paris")},
	},
	{
		spec: "TZ=Europe/Paris @sunset-1h10m 48.856614 2.3522219",
		want: &Schedule{Event: Sunset, Offset: -(time.Hour + 10*time.Minute), Lat: 48.856614, Lon: 2.3522219, Location: mustLoc("Europe/Paris")},
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		got, err := Parser{}.Parse(test.spec)
		if err != nil {
			if fmt.Sprint(err) != fmt.Sprint(test.wantErr) {
				t.Errorf("unexpected error parsing %q: got:%v want:%v",
					test.spec, err, test.wantErr)
			}
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("unexpected result for %q:\ngot: %#v\nwant:%#v",
				test.spec, got, test.want)
		}
	}
}

var nextTests = []struct {
	name string
	spec string
	time time.Time
	want time.Time
}{
	{
		name: "before event",
		spec: "@sunrise 48.856614 2.3522219",
		time: time.Date(2021, 7, 30, 6, 19, 0, 0, mustLoc("Europe/Paris")),
		want: time.Date(2021, 7, 30, 6, 20, 5, 0, mustLoc("Europe/Paris")),
	},
	{
		name: "at event",
		spec: "@sunrise 48.856614 2.3522219",
		time: time.Date(2021, 7, 30, 6, 20, 5, 0, mustLoc("Europe/Paris")),
		want: time.Date(2021, 7, 31, 6, 21, 25, 0, mustLoc("Europe/Paris")),
	},
	{
		name: "after event",
		spec: "@sunrise 48.856614 2.3522219",
		time: time.Date(2021, 7, 30, 6, 21, 0, 0, mustLoc("Europe/Paris")),
		want: time.Date(2021, 7, 31, 6, 21, 25, 0, mustLoc("Europe/Paris")),
	},
}

func TestNext(t *testing.T) {
	for _, test := range nextTests {
		s, err := Parser{}.Parse(test.spec)
		if err != nil {
			t.Errorf("unexpected error parsing %q: got:%v", test.spec, err)
			continue
		}
		got := s.Next(test.time)
		if !got.Equal(test.want) {
			t.Errorf("unexpected next event for %s %q at %v: got:%v want:%v",
				test.name, test.spec, test.time, got, test.want)
		}
	}
}
