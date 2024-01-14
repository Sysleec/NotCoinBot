package main

type cliCommands struct {
	name        string
	description string
	callback    func(cfg *Config, args ...string) error
}

func commands() map[string]cliCommands {
	return map[string]cliCommands{
		"help": {
			name:        "help",
			description: "Display this help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the NotCoinBot",
			callback:    commandExit,
		},
		"start": {
			name:        "start",
			description: "Start the NotCoinBot",
			callback:    commandStart,
		},
		"create": {
			name:        "create <name>",
			description: "Creating a tg session (name can be any)",
			callback:    commandCreate,
		},
	}
}
