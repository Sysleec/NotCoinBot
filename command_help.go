package main

import "fmt"

func commandHelp(cfg *Config, args ...string) error {
	fmt.Println("Welcome to NotCoinBot!")
	fmt.Println()
	fmt.Println("Available commands:")
	for _, cmd := range commands() {
		fmt.Printf("  %s - %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}
