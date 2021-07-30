// Copyright ©2021 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sun_test

import (
	"fmt"
	"log"
	"time"

	"github.com/kortschak/sun"
)

func ExampleTimes() {
	// Calculate solar times for Paris on 2021-07-30.
	loc, err := time.LoadLocation("CET")
	if err != nil {
		log.Fatal(err)
	}
	date := time.Date(2021, 7, 30, 12, 0, 0, 0, loc)
	rise, noon, set := sun.Times(date, 48.856614, 2.3522219)

	fmt.Printf("Sunrise: %v\nNoon:    %v\nSunset:  %v\n", rise, noon, set)

	// Output:
	//
	// Sunrise: 2021-07-30 06:20:05 +0200 CEST
	// Noon:    2021-07-30 13:57:09 +0200 CEST
	// Sunset:  2021-07-30 21:34:13 +0200 CEST
}
