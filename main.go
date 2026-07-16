package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
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

var caughtPoks map[string]Pokemon

var loc Locations

var poks Pokemonsters

var pok Pokemon

const baseTime = 50000 * time.Millisecond

const baseLocation = "https://pokeapi.co/api/v2/location-area/"

const k = 2.0

var cache = pokecache.NewCache(baseTime)

func main() {
	caughtPoks = make(map[string]Pokemon)
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
		"catch": {
			name:        "catch",
			description: "catching pokemon",
			callback:    catch,
		},
		"inspect": {
			name:        "inspect",
			description: "displaying pokemons stats",
			callback:    inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "displaying caught pokemons",
			callback:    pokedex,
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
			body,_ = getMap(baseLocation)
		}
		loc = unmarshalJSON[Locations](body)
		cache.Add(baseLocation, body)
	} else {
		body, ok := cache.Get(config.Next)
		if !ok {
			body,_ = getMap(config.Next)
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
			body,_ = getMap(config.Previous)
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
	link := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", userInput)
	body, ok := cache.Get(link)

	if !ok {
		body,_ = getMap(link)
	}
	poks = unmarshalJSON[Pokemonsters](body)
	cache.Add(link, body)

	for _, val := range poks.Poks {
		fmt.Printf("%s\n", val.Pokemon.Name)
	}
	return nil
}

func catch(config *Config, userInput string) error {
	link := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", userInput)
	body, ok := cache.Get(link)

	if !ok {
		body, _ = getMap(link)
	}
	pok = unmarshalJSON[Pokemon](body)
	cache.Add(link, body)
	fmt.Printf("Throwing a Pokeball at %s...\n", userInput)
	base_exp := pok.BaseExperience
	roll := rand.Intn(101)
	difficulty := 1.0 + 99.0 * (math.Exp(k * (float64(base_exp) / 608.0)) - 1.0) / (math.Exp(k) - 1.0)
	if roll <= int(difficulty) {
		fmt.Printf("%s escaped!\n", userInput)
	} else {
		fmt.Printf("%s was caught!\n", userInput)
		caughtPoks[pok.Forms[0].Name] = pok
	}
	return nil
}

func inspect(config *Config, userInput string) error {
	if pok, ok := caughtPoks[userInput]; !ok {
		fmt.Printf("you have not caught that pokemon\n")
	}else{
		fmt.Printf("Name: %s\n", userInput)
		fmt.Printf("Height: %d\n", pok.Height)
		fmt.Printf("Weight: %d\n", pok.Weight)
		fmt.Printf("Stats:\n")
		for _, s := range pok.Stats {
    		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
		}
		fmt.Printf("Types: \n")
		for _, t := range pok.Types {
    		fmt.Printf(" - %s\n", t.Type.Name)
		}

	}
	return nil
}

func pokedex(config *Config, userInput string) error{
	if len(caughtPoks) > 0{
	for key, _ := range caughtPoks {
    		fmt.Printf("  -%s\n", key)
		}
	} else {
		fmt.Printf("No pokemons caught yet\n")
	}
		return nil
}

func getMap(link string) ([]byte, error) {
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Likely wrong name entered")
		return nil, err
	}
	if err != nil {
		fmt.Printf("Something went wrong")
		return nil, err
	}

	return body, nil
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

type Pokemon struct {
	BaseExperience int `json:"base_experience"`

	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`

	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []struct {

		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
