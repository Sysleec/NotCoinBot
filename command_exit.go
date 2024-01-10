package main

import (
	"fmt"
	"os"
)

func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Closing NotCoinBot...")
	os.Exit(0)
	return nil
}
