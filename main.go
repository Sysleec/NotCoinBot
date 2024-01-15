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
	fmt.Println("For start you need to create sessions with 'create <name>'")
	fmt.Println("For help write 'help'")
	color.Unset()
	cfg := Config{}

	Repl(&cfg)
}
