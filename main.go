package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

var commands map[string]cliCommand

var loc Locations

func main() {
	
	config := Config{}
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "displaying all commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displaying 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displaying previous 20 location areas",
			callback:    commandMapB,
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
		if len(userCleanInput) != 0 {
			if command, ok := commands[userCleanInput[0]]; ok {
				command.callback(&config)
			}
		}

	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Printf("\n")
	for key, val := range commands {
		fmt.Printf("%s: %s\n", key, val.description)
	}
	return nil
}

func commandMap(config *Config) error {

	if config.Next == "" {
		loc = getMap("https://pokeapi.co/api/v2/location-area/")
	}else{
		loc = getMap(config.Next)
	}
	config.Next = loc.Next
	config.Previous = loc.Previous
	for _, val := range loc.Results {
		fmt.Printf("%s\n", val["name"])
	}
	return nil
}

func commandMapB(config *Config) error {

	if config.Previous == "" {
		fmt.Printf("You are on the first page\n")
		return nil
	} else {
		loc = getMap(config.Previous)
	}
	config.Next = loc.Next
	config.Previous = loc.Previous
	for _, val := range loc.Results {
		fmt.Printf("%s\n", val["name"])
	}
	return nil
}

func getMap(link string) Locations {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	loc := Locations{}
	err = json.Unmarshal(body, &loc)
	if err != nil {
		fmt.Println(err)
	}
	return loc
}

type Locations struct {
	Count    int                 `json:"count"`
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Results  []map[string]string `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}
