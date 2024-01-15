package main

import (
	"fmt"

	"github.com/fatih/color"
)

func commandHelp(cfg *Config, args ...string) error {
	color.Set(color.FgHiYellow)
	fmt.Println("Welcome to NotCoinBot!")
	fmt.Println()
	fmt.Println("Available commands:")
	for _, cmd := range commands() {
		fmt.Printf("  %s - %s\n", cmd.name, cmd.description)
	}
	color.Unset()
	fmt.Println()
	return nil
}
