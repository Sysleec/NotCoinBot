package main

import (
	"fmt"

	"github.com/fatih/color"
)

type Config struct {
}

func main() {
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
	color.Unset()
	cfg := Config{}

	Repl(&cfg)
}
