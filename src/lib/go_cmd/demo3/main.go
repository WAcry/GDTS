package main

import (
	"context"
	"fmt"
	"lib/config"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

// What is happening in this lib:
// execute a command in a sub goroutine, which echo hello after 2 seconds
// kill that sub goroutine after 1 second, so it will not print hello
func main() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		cmd    *exec.Cmd
		c      chan *result
		res    *result
	)

	// create a result channel
	c = make(chan *result, 1000)

	// context:   chan byte
	// cancel:  close(chan byte)

	ctx, cancel = context.WithCancel(context.TODO())

	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, config.BASH_PATH, "-c", "/usr/bin/sleep 2;echo hello;")

		// execute command, capture the output
		output, err = cmd.CombinedOutput()

		// send result to main goroutine
		c <- &result{
			err:    err,
			output: output,
		}
	}()

	// sleep 1 second
	time.Sleep(1 * time.Second)

	// cancel context
	cancel()

	// in main goroutine, wait for sub goroutine to exit, and print job execution result
	res = <-c

	// print job execute result
	fmt.Println(res.err, string(res.output))
}
