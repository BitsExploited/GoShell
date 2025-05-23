package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"bytes"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func isBuiltIn(command string) bool {
	builtIns := []string{"exit", "echo", "type", "pwd", "cd"}
	for _, i := range builtIns {
		if command == i {
			return true
		}
	}
	return false
}

func executeCommand(name string, args []string) {
	_, found := checkFileInPath(name)
	if found {
		var outb bytes.Buffer
		cmd := exec.Command(name, args...)
		cmd.Stdout = &outb
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s: could not execute process\n", name)
			return
		}
		fmt.Print(outb.String())
	} else {
		fmt.Printf("%s: command not found\n", name)
	}
}

func checkFileInPath(file string) (string, bool) {
	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, pathDir := range paths {
		if _, err := os.Stat(filepath.Join(pathDir, file)); err == nil {
			return pathDir, true
		}
	}
	return "", false
}

func exitCommand(args []string) {
	if len(args) > 0 {
		value, err := strconv.Atoi(args[0])
		if err != nil {
			os.Exit(0)
		}
		os.Exit(value)
	}
}

func cdCommand(args []string) {
	if len(args) == 0 || args[0] == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("error getting home directory: %v\n", err)
			return
		}
		err = os.Chdir(home)
		if err != nil {
			fmt.Printf("error changing to home directory: %v\n", err)
		}
		return
	}
	err := os.Chdir(args[0])
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", args[0])
		return
	}
}

func pwdCommand() {
	pwd , err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(pwd)
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
	return
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

		switch command {
		case "cd":
			cdCommand(args)
		case "pwd":
			pwdCommand()
		case "exit":
			exitCommand(args)
		case "echo":
			echoCommand(args)
		case "type":
			typeCommand(args[0], path)
		default:
			executeCommand(command, args)
		}
	}
}
