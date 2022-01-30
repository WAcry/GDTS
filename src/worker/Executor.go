package worker

import (
	"common"
	"math/rand"
	"os/exec"
	"runtime"
	"time"
)

type WkExecutor struct {
}

var (
	Executor *WkExecutor
)

func (executor *WkExecutor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)

		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		// init job lock
		jobLock = JobManager.CreateJobLock(info.Job.Name)

		// record job start time
		result.StartTime = time.Now()

		// lock with random sleep (because some machine may be a little earlier than others, so a random delay helps to split the works more equally)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		err = jobLock.TryLock()
		defer jobLock.Unlock()

		if err != nil { // lock failed
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// lock success, reset job start time
			result.StartTime = time.Now()

			// shell command

			sysType := runtime.GOOS
			if sysType == "windows" {
				cmd = exec.CommandContext(info.CancelCtx, Config.WinBashPath, "-c", info.Job.Command)
			} else {
				cmd = exec.CommandContext(info.CancelCtx, Config.BashPath, "-c", info.Job.Command)
			}

			// execute command and get output
			output, err = cmd.CombinedOutput()

			// record job end time
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		// after job finished, send result to Scheduler. Scheduler will delete the executing record from executingTable
		Scheduler.PushJobResult(result)
	}()
}

func InitExecutor() (err error) {
	Executor = &WkExecutor{}
	return
}
