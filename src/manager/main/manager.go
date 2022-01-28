package main

import (
	"flag"
	"fmt"
	"manager"
	"runtime"
	"time"
)

var (
	confFile string // Config file path
)

// parse command line arguments and initialize thread amount
func initManager() {
	// manager -config ./manager.json
	// manager -h
	// default path is set as src/manager/main/manager.json here
	flag.StringVar(&confFile, "config", "src/manager/main/manager.json", "path of manager.json")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	initManager()

	// load config
	if err = manager.InitConfig(confFile); err != nil {
		goto ERR
	}

	// initialize worker discovery module
	if err = manager.InitWorkerManager(); err != nil {
		goto ERR
	}

	// initialize logger module
	if err = manager.InitLogManager(); err != nil {
		goto ERR
	}

	//  initialize job module
	if err = manager.InitJobManager(); err != nil {
		goto ERR
	}

	// start api controller
	if err = manager.InitController(); err != nil {
		goto ERR
	}

	// loop forever
	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)
}
