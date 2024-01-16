package main

import (
	"fmt"
	"github.com/Sysleec/NotCoinBot/internal/clicker"
)

func commandStart(cfg *Config, args ...string) error {
	fmt.Println("Starting bots...")
	clicker.ClickerStart()
	return nil
}
