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
	"time"

	"github.com/skookum-ua/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

var commands map[string]cliCommand

var loc Locations

var poks Pokemonsters

const baseTime = 50000 * time.Millisecond

const baseLocation = "https://pokeapi.co/api/v2/location-area/"

var cache = pokecache.NewCache(baseTime)

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
		"explore": {
			name:        "explore",
			description: "displaying pokemons in area",
			callback:    explore,
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
			arg := ""
			if len(userCleanInput) > 1 {
				arg = userCleanInput[1]
			}
			if command, ok := commands[userCleanInput[0]]; ok {
				command.callback(&config, arg)
			}
		}

	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(config *Config, userInput string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, userInput string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Printf("\n")
	for key, val := range commands {
		fmt.Printf("%s: %s\n", key, val.description)
	}
	return nil
}

func commandMap(config *Config, userInput string) error {


	if config.Next == "" {
		body, ok := cache.Get(baseLocation)
		if !ok {
			body = getMap(baseLocation)
		}
		loc = unmarshalJSON[Locations](body)
		cache.Add(baseLocation, body)
	} else {
		body, ok := cache.Get(config.Next)
		if !ok {
			body = getMap(config.Next)
		}
		loc = unmarshalJSON[Locations](body)
		cache.Add(config.Next, body)
	}
	config.Next = loc.Next
	config.Previous = loc.Previous
	for _, val := range loc.Results {
		fmt.Printf("%s\n", val["name"])
	}
	return nil
}

func commandMapB(config *Config, userInput string) error {

	if config.Previous == "" {
		fmt.Printf("You are on the first page\n")
		return nil
	} else {
		body, ok := cache.Get(config.Previous)
		if !ok {
			body = getMap(config.Previous)
		}
		loc = unmarshalJSON[Locations](body)
		cache.Add(config.Previous, body)
	}
	config.Next = loc.Next
	config.Previous = loc.Previous
	for _, val := range loc.Results {
		fmt.Printf("%s\n", val["name"])
	}
	return nil
}

func explore(config *Config, userInput string) error {
		link := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s",userInput)
		body, ok := cache.Get(link)
		fmt.Printf("%s\n",link)
		if !ok {
			body = getMap(link)
		}
		poks = unmarshalJSON[Pokemonsters](body)
		cache.Add(link, body)

	for _, val := range poks.Poks {
		fmt.Printf("%s\n", val.Pokemon.Name)
	}
	return nil
}

func getMap(link string) []byte {
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
	
	return body
}

func unmarshalJSON[T any](body []byte) T {
	var result T
	err := json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Помилка парсингу JSON:", err)
	}
	return result
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

type Pokemonsters struct {
	Poks []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}