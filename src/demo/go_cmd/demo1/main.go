package main

import (
	"demo/config"
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
	)

	cmd = exec.Command(config.BASH_PATH, "-c", "echo 1")

	err = cmd.Run()

	fmt.Println(err)
}
