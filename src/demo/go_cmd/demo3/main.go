package main

import (
	"GDTS/src/config"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

// What is happening in this demo:
// execute a command in a sub goroutine, which echo hello after 2 seconds
// kill that sub goroutine after 1 second, so it will not print hello
func main() {
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		cmd        *exec.Cmd
		resultChan chan *result
		res        *result
	)

	// create a result channel
	resultChan = make(chan *result, 1000)

	// context:   chan byte
	// cancelFunc:  close(chan byte)

	ctx, cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, config.Path["bash"], "-c", "/usr/bin/sleep 2;echo hello;")

		// execute command, capture the output
		output, err = cmd.CombinedOutput()

		// send result to main goroutine
		resultChan <- &result{
			err:    err,
			output: output,
		}
	}()

	// sleep 1 second
	time.Sleep(1 * time.Second)

	// cancel context
	cancelFunc()

	// in main goroutine, wait for sub goroutine to exit, and print task execution result
	res = <-resultChan

	// print task execute result
	fmt.Println(res.err, string(res.output))
}
