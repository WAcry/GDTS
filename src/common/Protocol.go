package common

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"utils/gorhill/cronexpr"
)

// Job Timed Task
type Job struct {
	Name     string `json:"name"`     //  job name
	Command  string `json:"command"`  // shell command
	CronExpr string `json:"cronExpr"` // cron schedule expression
}

type JobSchedulePlan struct {
	Job      *Job                 // job info
	Expr     *cronexpr.Expression // cron expression
	NextTime time.Time            // next execution time
}

type JobExecuteInfo struct {
	Job        *Job               // job status
	PlanTime   time.Time          // planned time
	RealTime   time.Time          // real execution time
	CancelCtx  context.Context    // command context
	CancelFunc context.CancelFunc //  cancel function used to cancel command
}

// Response Http Response
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// JobEvent   SAVE = 1, DELETE = 2, KILL = 3
type JobEvent struct {
	EventType int
	Job       *Job
}

type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // execute status
	Output      []byte          // command output
	Err         error           // command error
	StartTime   time.Time       // start time
	EndTime     time.Time       // end time
}

type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`           // job name
	Command      string `json:"command" bson:"command"`           // command
	Err          string `json:"err" bson:"err"`                   // error
	Output       string `json:"output" bson:"output"`             // shell output
	PlanTime     int64  `json:"planTime" bson:"planTime"`         // planned time
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"` // real execution time
	StartTime    int64  `json:"startTime" bson:"startTime"`       // start time of execution
	EndTime      int64  `json:"endTime" bson:"endTime"`           // end time of execution
}

type LogBatch struct {
	Logs []interface{} // multiple logs
}

type JobLogFilter struct {
	JobName string `bson:"jobName"`
}

// SortLogByStartTime log sort rule by start time
type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"` // {startTime: -1}
}

func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	// serialize json
	resp, err = json.Marshal(response)
	return
}

// UnpackJob deserialize job json
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// ExtractJobName extract job name from etcd's key
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

// ExtractKillerName e.g extract "job_x" from /cron/killer/job_x
func ExtractKillerName(killerKey string) string {
	return strings.TrimPrefix(killerKey, JOB_KILLER_DIR)
}

// BuildJobEvent 1)update 2)delete
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var (
		expr *cronexpr.Expression
	)

	// parser cron expression
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	// generate job schedule plan object
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime, // planned time
		RealTime: time.Now(),               // real execution time
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}

func ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JOB_WORKER_DIR)
}
