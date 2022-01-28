package main

import (
	"fmt"
	"time"
	"utils/gorhill/cronexpr"
)

func main() {
	var (
		exp *cronexpr.Expression
		err error
		now time.Time
		nxt time.Time
	)

	// linux crontab
	// cron: [second 0-59] minute 0-59 hour 0-23 day 1-31 month 1-12 weekday 0-6 [year 1970-2099]
	// */5 * * * *		: every 5 minutes
	// 0 13 * * 0-2		: every Sunday to Tuesday at 13pm
	// 0 13,16 * * *	: every 13pm, 16pm
	// * * * * * * 2022 : every second of 2022

	// run every 5 seconds
	if exp, err = cronexpr.Parse("*/5 * * * * * *"); err != nil {
		fmt.Println(err)
		return
	}

	// at second 0, 5, 10, ... , 45, 50, 55

	// current time
	now = time.Now()
	// next time
	nxt = exp.Next(now)

	// wait for the timer to expire
	time.AfterFunc(nxt.Sub(now), func() { // run func after "nxt - now" seconds
		fmt.Println("run once and exit:", nxt)
	})

	time.Sleep(5 * time.Second)
}
