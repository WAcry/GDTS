package main

import (
	"fmt"
	"time"
	"utils/gorhill/cronexpr"
)

// Job represents a cron job
type Job struct {
	expr *cronexpr.Expression
	nxt  time.Time // expr.Next(now)
}

// need a scheduler goroutine to check all cron jobs: who is expired, execute it
func main() {

	var (
		job      *Job
		exp      *cronexpr.Expression
		now      time.Time
		name2job map[string]*Job // key: job name
	)

	name2job = make(map[string]*Job)

	// current time
	now = time.Now()

	// First, define two cron jobs
	exp = cronexpr.MustParse("*/3 * * * * * *")
	job = &Job{
		expr: exp,
		nxt:  exp.Next(now),
	}
	// register cron job
	name2job["job1"] = job

	exp = cronexpr.MustParse("*/5 * * * * * *")
	job = &Job{
		expr: exp,
		nxt:  exp.Next(now),
	}
	// register cron job
	name2job["job2"] = job

	// start a goroutine
	go func() {
		var (
			name string
			job  *Job
			now  time.Time
		)

		// check the schedule table every 100 milliseconds
		for {
			now = time.Now()

			for name, job = range name2job {
				// check if the job is expired
				if job.nxt.Before(now) || job.nxt.Equal(now) {
					// start a goroutine to execute the job
					go func(name string) { // declare inner func
						fmt.Println("execute:", name)
					}(name) // pass name to inner func

					// update the next time
					job.nxt = job.expr.Next(now)
					fmt.Println(name, "next execute time:", job.nxt)
				}
			}

			// sleep 100 milliseconds using channel // alternative: time.Sleep(100 * time.Millisecond)
			select { // when mutiple channels are ready, select will select one of them randomly
			case <-time.NewTimer(100 * time.Millisecond).C: // will be blocked for 100 milliseconds
				// 	case <-(make(chan int)) : // another channel, will be blocked until the channel is ready
			}
		}
	}()

	// sleep so the child goroutine can run for 60 seconds only
	time.Sleep(60 * time.Second)
}
