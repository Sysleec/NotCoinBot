package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Repl(cfg *Config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("NotCoinBot >")
		scanner.Scan()

		command := scanner.Text()
		if len(command) == 0 {
			continue
		}

		valCommand := validateCommand(command)
		commandName := valCommand[0]
		args := []string{}
		if len(valCommand) > 1 {
			args = valCommand[1:]
		}

		cmd, ok := commands()[commandName]
		if ok {
			err := cmd.callback(cfg, args...)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
			continue
		}

	}
}

func validateCommand(str string) []string {
	lowStr := strings.ToLower(str)
	return strings.Fields(lowStr)
}
