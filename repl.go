package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func Repl(cfg *Config) {
	color.Set(color.FgHiYellow)
	data := []byte(`
	           ████████████
	       ████░░░░░░░░░░░░████
	     ██░░░░▒▒▒▒▒▒▒▒▒▒▒▒░░░░██
	   ██░░▒▒░░░░      ░░░░░░  ░░██
	 ██░░▒▒░░░░░░  ░░░░▒▒░░░░░░  ░░██
	 ██░░▒▒░░      ░░░░░░  ░░░░  ░░██
	 ██░░▒▒░░  ░░░░░░░░░░░░▒▒░░  ░░██
	 ██░░▒▒░░  ░░░░░░░░░░░░▒▒░░  ░░██
	 ██░░▒▒░░░░▒▒░░░░░░▒▒▒▒▒▒░░  ░░██
	 ██░░▒▒░░░░░░  ░░░░▒▒░░░░░░  ░░██
	   ██░░▒▒░░░░░░▒▒▒▒▒▒░░░░  ░░██
	   ██░░░░  ░░░░░░░░░░░░  ░░░░██
	     ██░░░░            ░░░░██
	       ████░░░░░░░░░░░░████
	           ████████████
	`)
	fmt.Println(string(data))

	color.Set(color.FgHiGreen)
	fmt.Println("Welcome to the NotCoinBot!!!")
	fmt.Println()
	fmt.Println("This bot is designed to automate the process of collecting coins in the game NotCoin")
	fmt.Println("Author - github.com/Sysleec")
	fmt.Println()
	fmt.Println("You can support me with a donation to my TON wallet")
	fmt.Println("UQDPJTL7YyUk8Jm92Vymg2X4V-pKtNdarKMyv-r1oYrcfJKX")
	fmt.Println("And also join our community NotCoin Lions")
	fmt.Println("https://t.me/NotcoinLions")
	fmt.Println()
	fmt.Println("For start you need to create sessions with 'create <name>'")
	fmt.Println("To run the clicker, write 'start'")
	fmt.Println("For help write 'help'")
	fmt.Println()
	color.Unset()

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
