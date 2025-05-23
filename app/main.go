package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(pwd)
}

func echoCommand(args []string) {
	echoStr := strings.Join(args, " ")
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

// parseCommand tokenizes the input string, respecting single and double quotes
func parseCommand(input string) []string {
	var result []string
	var current strings.Builder
	var inSingleQuote, inDoubleQuote bool
	input = strings.TrimSpace(input)

	for i := 0; i < len(input); i++ {
		c := input[i]

		if inSingleQuote {
			if c == '\'' {
				inSingleQuote = false
			} else {
				current.WriteByte(c)
			}
			continue
		}

		if inDoubleQuote {
			if c == '"' {
				inDoubleQuote = false
			} else {
				current.WriteByte(c)
			}
			continue
		}

		if c == '\'' {
			inSingleQuote = true
			continue
		}

		if c == '"' {
			inDoubleQuote = true
			continue
		}

		if c == ' ' || c == '\t' {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(c)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')

		commandParts := parseCommand(command)
		if len(commandParts) == 0 {
			continue
		}

		path := os.Getenv("PATH")
		commandName := commandParts[0]
		args := commandParts[1:]

		switch commandName {
		case "cd":
			cdCommand(args)
		case "pwd":
			pwdCommand()
		case "exit":
			exitCommand(args)
		case "echo":
			echoCommand(args)
		case "type":
			if len(args) > 0 {
				typeCommand(args[0], path)
			} else {
				fmt.Println("type: missing argument")
			}
		default:
			executeCommand(commandName, args)
		}
	}
}
