package main

import (
	"bufio"
	"fmt"
	"os"
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

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading input: ", err)
			os.Exit(1)
		}
		command = strings.TrimSuffix(command, "\n")
		
		if command == "exit 0" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			phrase := command[5:]
			fmt.Println(phrase)
		} else if strings.HasPrefix(command, "type") {
			phrase := command[5:]
			if isBuiltIn(phrase) {
				fmt.Println(phrase + " is a shell builtin")
			} else {
				fmt.Printf("%s: not found\n", phrase)
			}
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}
