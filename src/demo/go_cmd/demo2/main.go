package main

import (
	"demo/config"
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
	cmd = exec.Command(config.BASH_PATH, "-c", "/usr/bin/sleep 5;/usr/bin/pwd;/usr/bin/ls -l")

	// execute command, capture the output of the child process from the pipe
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	// print output of the child process
	fmt.Println(string(output))

	cmd = exec.Command(config.PYTHON_PATH, config.HELLO_PY_PATH)

	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}
