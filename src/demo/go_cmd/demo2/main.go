package main

import (
	"GDTS/config"
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	// Generate cmd
	cmd = exec.Command(config.Path["bash"], "-c", "/usr/bin/sleep 5;/usr/bin/pwd;/usr/bin/ls -l")

	// execute command, capture the output of the child process from the pipe
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	// print output of the child process
	fmt.Println(string(output))

	cmd = exec.Command(config.Path["python"], config.Path["hello.py"])

	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}
