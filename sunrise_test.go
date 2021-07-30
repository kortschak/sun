// Copyright ©2021 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sun

import (
	"testing"
	"time"
)

var timesTests = []struct {
	name     string
	date     time.Time
	lat, lon float64

	// Want values obtained from https://www.timeanddate.com/.
	wantRise time.Time
	wantNoon time.Time
	wantSet  time.Time
}{
	{
		name: "Ochre Point July", lat: -35.2163355, lon: 138.4751278,
		date:     time.Date(2021, 7, 30, 12, 0, 0, 0, mustLoc("Australia/Adelaide")),
		wantRise: time.Date(2021, 7, 30, 7, 11, 0, 0, mustLoc("Australia/Adelaide")),
		wantNoon: time.Date(2021, 7, 30, 12, 22, 0, 0, mustLoc("Australia/Adelaide")),
		wantSet:  time.Date(2021, 7, 30, 17, 32, 0, 0, mustLoc("Australia/Adelaide")),
	},
	{
		name: "Ochre Point February", lat: -35.2163355, lon: 138.4751278,
		date:     time.Date(2021, 2, 21, 12, 0, 0, 0, mustLoc("Australia/Adelaide")),
		wantRise: time.Date(2021, 2, 21, 6, 55, 0, 0, mustLoc("Australia/Adelaide")),
		wantNoon: time.Date(2021, 2, 21, 13, 29, 0, 0, mustLoc("Australia/Adelaide")),
		wantSet:  time.Date(2021, 2, 21, 20, 2, 0, 0, mustLoc("Australia/Adelaide")),
	},
	{
		name: "New York July", lat: 40.6976637, lon: -74.119764,
		date:     time.Date(2021, 7, 30, 12, 0, 0, 0, mustLoc("America/New_York")),
		wantRise: time.Date(2021, 7, 30, 5, 51, 0, 0, mustLoc("America/New_York")),
		wantNoon: time.Date(2021, 7, 30, 13, 2, 0, 0, mustLoc("America/New_York")),
		wantSet:  time.Date(2021, 7, 30, 20, 13, 0, 0, mustLoc("America/New_York")),
	},
	{
		name: "New York February", lat: 40.6976637, lon: -74.119764,
		date:     time.Date(2021, 2, 21, 12, 0, 0, 0, mustLoc("America/New_York")),
		wantRise: time.Date(2021, 2, 21, 6, 41, 0, 0, mustLoc("America/New_York")),
		wantNoon: time.Date(2021, 2, 21, 12, 9, 0, 0, mustLoc("America/New_York")),
		wantSet:  time.Date(2021, 2, 21, 17, 38, 0, 0, mustLoc("America/New_York")),
	},
}

func TestTimes(t *testing.T) {
	const tol = 5 * time.Minute

	for _, test := range timesTests {
		gotRise, gotNoon, gotSet := Times(test.date, test.lat, test.lon)
		if !similarTime(gotRise, test.wantRise, tol) {
			t.Errorf("unexpected sun rise time for %q: got:%v want:%v", test.name, gotRise, test.wantRise)
		}
		if !similarTime(gotNoon, test.wantNoon, tol) {
			t.Errorf("unexpected solar for %q: got:%v want:%v", test.name, gotNoon, test.wantNoon)
		}
		if !similarTime(gotSet, test.wantSet, tol) {
			t.Errorf("unexpected sun set time for %q: got:%v want:%v", test.name, gotSet, test.wantSet)
		}
	}
}

func similarTime(a, b time.Time, tol time.Duration) bool {
	if a.Before(b) {
		return b.Sub(a) < tol
	}
	return a.Sub(b) < tol
}

func mustLoc(name string) *time.Location {
	// return time.Local
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
