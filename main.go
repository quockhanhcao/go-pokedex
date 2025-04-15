package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/quockhanhcao/go-pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	NextURL     string
	PreviousURL string
	Location    string
	Pokemon     string
}

var caughtPokemon map[string]pokeapi.Pokemon

var supportedCommand map[string]cliCommand

func init() {
	supportedCommand = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations in Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations in Pokemon world",
			callback:    commandMapback,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location in the Pokemon world",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon",
			callback:    commandInspect,
		},
	}
	caughtPokemon = map[string]pokeapi.Pokemon{}
}

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	pokeapi.CloseCache()
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, command := range supportedCommand {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *config) error {
	data, err := pokeapi.GetLocationAreaData(config.NextURL)
	if err != nil {
		return err
	}
	config.NextURL = data.Next
	if data.Previous != nil {
		config.PreviousURL = *data.Previous
	}
	for _, result := range data.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapback(config *config) error {
	data, err := pokeapi.GetLocationAreaData(config.PreviousURL)
	if err != nil {
		return err
	}
	config.NextURL = data.Next
	if data.Previous != nil {
		config.PreviousURL = *data.Previous
	}
	for _, result := range data.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(config *config) error {
	data, err := pokeapi.GetLocationDetailedData(config.Location)
	if err != nil {
		return err
	}
	for _, pokemonEncounter := range data.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config) error {
	data, err := pokeapi.GetPokemonStats(config.Pokemon)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)
	catchChance := 100 - int((float64(data.BaseExperience)/500.0)*100)
	randNum := rand.Intn(100) + 1
	if randNum <= catchChance {
		fmt.Printf("%s was caught!\n", data.Name)
		caughtPokemon[data.Name] = data
	} else {
		fmt.Printf("%s escaped!\n", data.Name)
	}
	return nil
}

func commandInspect(config *config) error {
	data, err := pokeapi.GetPokemonStats(config.Pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", data.Name)
	fmt.Printf("Height: %d\n", data.Height)
	fmt.Printf("Weight: %d\n", data.Weight)
	fmt.Println("Stats:")
	for _, stat := range data.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range data.Types {
		fmt.Printf("  -%s\n", t.Type.Name)
	}

	return nil
}

func main() {
	config := config{
		NextURL:     "https://pokeapi.co/api/v2/location-area",
		PreviousURL: "",
		Location:    "",
		Pokemon:     "",
	}
	defer pokeapi.CloseCache()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		data := scanner.Text()
		cleanedInput := cleanInput(data)
		if len(cleanedInput) > 0 {
			if command, exists := supportedCommand[cleanedInput[0]]; exists {
				if len(cleanedInput) > 1 {
					if cleanedInput[0] == "explore" {
						config.Location = cleanedInput[1]
					} else if cleanedInput[0] == "catch" || cleanedInput[0] == "inspect" {
						config.Pokemon = cleanedInput[1]
					}
				}
				command.callback(&config)
			} else {
				fmt.Println("Unknown command")
			}
		}
	}

}

func cleanInput(text string) []string {
	var cleaned []string
	trimmed := strings.TrimSpace(text)
	words := strings.Split(trimmed, " ")

	for _, word := range words {
		if word != "" {
			cleaned = append(cleaned, strings.ToLower(word))
		} else {
			fmt.Println("skipping empty word")
		}
	}
	return cleaned
}
