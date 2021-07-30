// Copyright ©2021 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sun provides solar time calculations.
package sun

import (
	"math"
	"time"
)

// Times returns the times for sun rise, solar noon and sun set for the given
// time and location based on the formulae provided by the NOAA.
//
// See https://gml.noaa.gov/grad/solcalc/solareqns.PDF
func Times(t time.Time, lat, lon float64) (rise, noon, set time.Time) {
	gamma := fractionalYear(t)
	decl := solarDeclinationAngle(gamma)
	ha := deg(hourAngle(lat, decl))
	eqTime := equationOfTime(gamma)
	rise = timeFromMinutes(t, 720-4*(lon+ha)-eqTime)
	noon = timeFromMinutes(t, 720-4*(lon)-eqTime)
	set = timeFromMinutes(t, 720-4*(lon-ha)-eqTime)
	return rise, noon, set
}

// fractionalYear returns the fractional year in radians.
func fractionalYear(t time.Time) float64 {
	hr, _, _ := t.Clock()
	return (2 * math.Pi / daysInYearFor(t)) * (float64(t.YearDay()-1) + float64(hr-12)/24)
}

// daysInYearFor returns the number of days in the time's year.
func daysInYearFor(t time.Time) float64 {
	year, _, _ := t.Date()
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return 366
	}
	return 365
}

// solarDeclinationAngle returns the solar declination angle in radians.
func solarDeclinationAngle(gamma float64) float64 {
	return 0.006918 - 0.399912*math.Cos(gamma) + 0.070257*math.Sin(gamma) - 0.006758*math.Cos(2*gamma) + 0.000907*math.Sin(2*gamma) - 0.002697*math.Cos(3*gamma) + 0.00148*math.Sin(3*gamma)
}

// equationOfTime returns the equation of time in minutes.
func equationOfTime(gamma float64) float64 {
	return 229.18 * (0.000075 + 0.001868*math.Cos(gamma) - 0.032077*math.Sin(gamma) - 0.014615*math.Cos(2*gamma) - 0.040849*math.Sin(2*gamma))
}

// hourAngle returns the hour angle in radians.
func hourAngle(lat, decl float64) float64 {
	latRad := rad(lat)
	return math.Acos((math.Cos(rad(90.833)) / (math.Cos(latRad) * math.Cos(decl))) - math.Tan(latRad)*math.Tan(decl))
}

// deg returns a radians as degrees.
func deg(a float64) float64 {
	return a / math.Pi * 180
}

// rad returns a degrees as radians.
func rad(a float64) float64 {
	return a / 180 * math.Pi
}

// timeFromMinutes returns the time corresponding to the minute offset into
// the given date.
func timeFromMinutes(date time.Time, minutes float64) time.Time {
	hour := int(minutes) / 60
	min := int(minutes) % 60
	_, s := math.Modf(minutes)
	sec := int(s * 60)
	u := date.In(time.UTC)
	return time.Date(u.Year(), u.Month(), u.Day(), hour, min, sec, 0, time.UTC).In(date.Location())
}
