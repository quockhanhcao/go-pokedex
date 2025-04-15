package main

import (
	"bufio"
	"fmt"
	"github.com/quockhanhcao/go-pokedex/internal/pokeapi"
	"os"
	"strings"
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
}

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
	}
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

func main() {
	config := config{
		NextURL:     "https://pokeapi.co/api/v2/location-area",
		PreviousURL: "",
		Location:    "",
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
					config.Location = cleanedInput[1]
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
