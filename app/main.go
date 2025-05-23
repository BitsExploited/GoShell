package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func isBuiltIn(command string) bool {
	builtIns := []string{"exit", "echo", "type"}
	for _, i := range builtIns {
		if command == i {
			return true
		}
	}
	return false
}

func echoCommand(argument []string) {
	echoStr := strings.Join(argument, " ")
	fmt.Println(echoStr)
}

func typeCommand(argument string, path string) {
	if isBuiltIn(argument) {
		fmt.Printf("%s is a shell builtin\n", argument)
		return
	}

	os.Setenv("PATH", path)
	if cmdPath, err := exec.LookPath(argument); err == nil {
		fmt.Printf("%s is %s\n", argument, cmdPath)
	} else {
		fmt.Printf("%s: not found\n", argument)
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		
		command = strings.TrimSpace(command)
		commandParts := strings.Fields(command)

		path := os.Getenv("PATH")
		command = commandParts[0]
		args := commandParts[1:]

		if command == "exit" {
			break

		} else if strings.HasPrefix(command, "echo") {
			echoCommand(args)

		} else if strings.HasPrefix(command, "type") {
			typeCommand(args[0], path)

		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
