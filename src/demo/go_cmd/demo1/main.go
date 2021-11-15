package main

import (
	"GDTS/config"
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
	)

	cmd = exec.Command(config.Path["bash"], "-c", "echo 1")

	err = cmd.Run()

	fmt.Println(err)
}
