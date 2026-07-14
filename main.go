package main
import (
	"fmt"
	"strings"
	"os"
	"bufio"
)

type cliCommand struct{
	name string
	description string
	callback func() error
}
var commands map[string]cliCommand

func main(){
	commands = map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "displaying all commands",
			callback: commandHelp,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			break
		}
		userInput := scanner.Text()
		userCleanInput := cleanInput(userInput)
		if len(userCleanInput) != 0{
		if command, ok := commands[userCleanInput[0]]; ok {
			command.callback()
		}
	}
		
	}
}

	func cleanInput(text string) []string {
		return strings.Fields(strings.ToLower(text))
	}

	func commandExit() error {
		fmt.Println("Closing the Pokedex... Goodbye!")
		os.Exit(0)
		return nil
	}

	func commandHelp()error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Println("Usage:")
			fmt.Printf("\n")
		for key, val := range commands{
			fmt.Printf("%s: %s\n", key, val.description)
		}
		return nil
	}

