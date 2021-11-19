package main

import (
	"GDTS/utils/gorhill/cronexpr"
	"fmt"
	"time"
)

// CronJob represents a cron job
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time // expr.Next(now)
}

// need a scheduler goroutine to check all cron jobs: who is expired, execute it
func main() {

	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob // key: task name
	)

	scheduleTable = make(map[string]*CronJob)

	// current time
	now = time.Now()

	// First, define two cron jobs
	expr = cronexpr.MustParse("*/3 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// register cron job
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// register cron job
	scheduleTable["job2"] = cronJob

	// start a goroutine
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)

		// check the schedule table every 100 milliseconds
		for {
			now = time.Now()

			for jobName, cronJob = range scheduleTable {
				// check if the job is expired
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// start a goroutine to execute the job
					go func(jobName string) { // declare inner func
						fmt.Println("execute:", jobName)
					}(jobName) // pass jobName to inner func

					// update the next time
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "next execute time:", cronJob.nextTime)
				}
			}

			// sleep 100 milliseconds using channel // alternative: time.Sleep(100 * time.Millisecond)
			select { // when mutiple channels are ready, select will select one of them randomly
			case <-time.NewTimer(100 * time.Millisecond).C: // will be blocked for 100 milliseconds
				// 	case <-(make(chan int)) : // another channel, will be blocked until the channel is ready
			}
		}
	}()
}
