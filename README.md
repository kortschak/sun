# Package sun implements solar time calculations.

[![Go Documentation](https://pkg.go.dev/badge/github.com/kortschak/sun.svg)](https://pkg.go.dev/github.com/kortschak/sun) [![Build status](https://github.com/kortschak/sun/workflows/Test/badge.svg)](https://github.com/kortschak/sun/actions)

The package provides a single function that returns the sunrise, sunset and solar noon for a given time and location, and a cron spec parser and cron scheduler type that allows cron jobs to be scheduled relative to solar-time events using [github.com/robfig/cron/v3](https://github.com/robfig/cron).

```
package main

import (
	"fmt"
	"log"

	"github.com/kortschak/sun"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithParser(sun.Parser{}))

	// Set a reminder to go for a walk each evening.
	_, err := c.AddFunc("@sunset-30m 48.856614 2.3522219", func() {
		fmt.Println("Take a walk along the Seine before sunset.")
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	select {}
}
```
