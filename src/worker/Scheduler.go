package worker

import (
	"common"
	"fmt"
	"time"
)

type WkScheduler struct {
	jobEventChan      chan *common.JobEvent              // etcd task event channel
	jobPlanTable      map[string]*common.JobSchedulePlan // task schedule planning table
	jobExecutingTable map[string]*common.JobExecuteInfo  // task execution table
	jobResultChan     chan *common.JobExecuteResult      // task result queue
}

var (
	Scheduler *WkScheduler
)

func (scheduler *WkScheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo  *common.JobExecuteInfo
		jobExecuting    bool
		jobExisted      bool
		err             error
	)
	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: // save job event
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE: // delete job event
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL: // kill job event
		// if job is running, kill the job
		if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecuting {
			jobExecuteInfo.CancelFunc()
		}
	}
}

// TryStartJob try execute job
func (scheduler *WkScheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	// schedule is different from execute
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)

	// if a task runs for 1 minute but schedule every second, it will only run once to avoid concurrent issue

	// if job is running, skip this job scheduling
	if jobExecuteInfo, jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		// fmt.Println("尚未退出,跳过执行:", jobPlan.Job.Name)
		return
	}

	// generate job execute info
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	// save executing status
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	fmt.Println("run job:", jobExecuteInfo.Job.Name, jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
	Executor.ExecuteJob(jobExecuteInfo)
}

// TrySchedule re-calculate job schedule status
func (scheduler *WkScheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)

	// if no job, sleep for a while
	if len(scheduler.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	// current time
	now = time.Now()

	for _, jobPlan = range scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) // update next time
		}

		// calculate the nearest due time
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	// next schedule time interval (nearest job schedule time - current time)
	scheduleAfter = (*nearTime).Sub(now)
	return
}

// deal with job result
func (scheduler *WkScheduler) handleJobResult(result *common.JobExecuteResult) {
	var (
		jobLog *common.JobLog
	)
	// delete executing status
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)

	// generate log
	if result.Err != common.ERR_LOCK_ALREADY_ACQUIRED {
		jobLog = &common.JobLog{
			JobName:      result.ExecuteInfo.Job.Name,
			Command:      result.ExecuteInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
	}

	fmt.Println("job is done:", result.ExecuteInfo.Job.Name, string(result.Output), result.Err)
}

func (scheduler *WkScheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *common.JobExecuteResult
	)

	// init once (1 second)
	scheduleAfter = scheduler.TrySchedule()

	// timer of schedule
	scheduleTimer = time.NewTimer(scheduleAfter)

	// timed task common.Job
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: // listen job change event
			// CRUD to jobList in memory
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: // recent job is due
		case jobResult = <-scheduler.jobResultChan: // listen to job result
			scheduler.handleJobResult(jobResult)
		}
		// schedule the job once
		scheduleAfter = scheduler.TrySchedule()
		// reset schedule interval
		scheduleTimer.Reset(scheduleAfter)
	}
}

func (scheduler *WkScheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

func InitScheduler() (err error) {
	Scheduler = &WkScheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}
	go Scheduler.scheduleLoop()
	return
}

// PushJobResult pass back the execution result
func (scheduler *WkScheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}
