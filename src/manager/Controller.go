package manager

import (
	"common"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)

type MstController struct {
	httpServer *http.Server
}

var (
	// Controller singleton controller
	Controller *MstController
)

// POST job={"name": "job1", "command": "echo hello", "cronExpr": "* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	// 1, parse the post data
	if err = req.ParseForm(); err != nil {
		goto ERR
	}
	// 2, get job field
	postJob = req.PostForm.Get("job")
	// 3, decode job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	// 4, save to etcd
	if oldJob, err = JobManager.SaveJob(&job); err != nil {
		goto ERR
	}
	// 5, return normal response ({"errno": 0, "msg": "", "data": {....}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	// 6, return error response ({"errno": -1, "msg": "", "data": {....}})
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// POST /job/delete   name=job1
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error // interface{}
		name   string
		oldJob *common.Job
		bytes  []byte
	)

	// POST:   a=1&b=2&c=3
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.PostForm.Get("name")

	if oldJob, err = JobManager.DeleteJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// list all available jobs
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*common.Job
		bytes   []byte
		err     error
	)

	// get job list
	if jobList, err = JobManager.ListJobs(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// POST /job/kill  name=job1
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)

	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	name = req.PostForm.Get("name")

	if err = JobManager.KillJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// query job logs
func handleJobLog(resp http.ResponseWriter, req *http.Request) {
	var (
		err        error
		name       string // job name
		skipParam  string // start from which log
		limitParam string // return how many logs
		skip       int
		limit      int
		logArr     []*common.JobLog
		bytes      []byte
	)

	// parse GET params
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// parse : GET /job/log?name=job10&skip=0&limit=10
	name = req.Form.Get("name")
	skipParam = req.Form.Get("skip")
	limitParam = req.Form.Get("limit")
	if skip, err = strconv.Atoi(skipParam); err != nil {
		skip = 0
	}
	if limit, err = strconv.Atoi(limitParam); err != nil {
		limit = 20
	}

	if logArr, err = LogManager.ListLog(name, skip, limit); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", logArr); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// get list of healthy workers
func handleWorkerList(resp http.ResponseWriter, req *http.Request) {
	var (
		workerArr []string
		err       error
		bytes     []byte
	)

	if workerArr, err = WorkerManager.ListWorkers(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", workerArr); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func InitController() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir // root dir for static files
		staticHandler http.Handler
	)

	// set router
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/worker/list", handleWorkerList)

	//  /index.html

	staticDir = http.Dir(Config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler)) //   ./webroot/index.html

	// start to listen
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(Config.ApiPort)); err != nil {
		return
	}

	// create http server
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(Config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(Config.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	Controller = &MstController{
		httpServer: httpServer,
	}

	// start http server
	go httpServer.Serve(listener)

	return
}
