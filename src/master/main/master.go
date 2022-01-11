package main

import (
	"flag"
	"fmt"
	"master"
	"runtime"
	"time"
)

var (
	confFile string // Config file path
)

// parse command line arguments and initialize thread amount
func initMaster() {
	// master -config ./master.json
	// master -h
	// default path is set as src/master/main/master.json here
	flag.StringVar(&confFile, "config", "src/master/main/master.json", "path of master.json")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	initMaster()

	// load config
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// initialize worker discovery module
	if err = master.InitWorkerManager(); err != nil {
		goto ERR
	}

	// initialize logger module
	if err = master.InitLogManager(); err != nil {
		goto ERR
	}

	//  initialize job module
	if err = master.InitJobManager(); err != nil {
		goto ERR
	}

	// start api controller
	if err = master.InitController(); err != nil {
		goto ERR
	}

	// loop forever
	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)
}
