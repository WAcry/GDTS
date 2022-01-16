package main

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)
	sysType := runtime.GOOS
	if sysType == "windows" {
		cmd = exec.CommandContext(context.TODO(), "C:\\cygwin64\\bin\\bash.exe", "-c", "/usr/bin/sleep 2;echo hello;")
	} else {
		exec.CommandContext(context.TODO(), "/bin/bash", "-c", "echo 1")
	}
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(output))
}
