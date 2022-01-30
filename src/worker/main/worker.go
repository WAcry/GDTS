package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
	"worker"
)

var (
	configPath string // config file path
)

// parse command line arguments and initialize thread amount
func initWorker() {
	// worker -config ./worker.json
	// worker -h
	// default path is set as src/worker/main/worker.json here
	flag.StringVar(&configPath, "config", "src/worker/main/worker.json", "worker.json")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	initWorker()

	// init config module
	if err = worker.InitConfig(configPath); err != nil {
		goto ERR
	}

	// start job executor module
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}

	// start job scheduler module
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}

	// start job manager module
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)
}
