package main

import (
	"fmt"

	"bufio"

	"os"

	"strings"

	"internal/pokeapi"

	"math/rand"
)

var pokedex map[string]pokeapi.Pokemon = map[string]pokeapi.Pokemon{}

type config struct {
	Next     string
	Previous string
	Params   []string
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	commands := getSupportedCommands()

	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(c *config) error {
	if c.Next == "" {
		fmt.Println("No more locations to display")
		return nil
	}

	data, err := pokeapi.GetLocationAreas(c.Next)
	if err != nil {
		fmt.Println(err)
		return err
	}

	c.Next = data.Next
	c.Previous = data.Previous

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack(c *config) error {
	if c.Previous == "" {
		fmt.Println("No more locations to display")
		return nil
	}

	data, err := pokeapi.GetLocationAreas(c.Previous)
	if err != nil {
		fmt.Println(err)
		return err
	}

	c.Next = data.Next
	c.Previous = data.Previous

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(c *config) error {
	if len(c.Params) == 0 {
		fmt.Println("Please provide a location name")
		return nil
	}

	locationName := c.Params[0]

	data, err := pokeapi.GetPokemonInLocationArea(locationName)

	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(data.PokemonEncounters) == 0 {
		fmt.Println("No pokemon found in this location")
		return nil
	}

	fmt.Println("Found Pokemon:")

	for _, pokemon := range data.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config) error {
	if len(c.Params) == 0 {
		fmt.Println("Please provide a pokemon name")
		return nil
	}

	pokemonName := c.Params[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := pokeapi.GetPokemonInfo(pokemonName)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	chance := rand.Intn(pokemon.BaseExperience)

	if chance > 50 {
		fmt.Printf("%s escaped!\n", pokemonName)
	} else {
		fmt.Printf("%s was caught!\n", pokemonName)
		pokedex[pokemon.Name] = pokemon
		fmt.Println("You may now inspect it with the inspect command.")
	}

	return nil
}

func commandInspect(c *config) error {
	if len(c.Params) == 0 {
		fmt.Println("Please provide a pokemon name")
		return nil
	}

	pokemonName := c.Params[0]

	pokemon, ok := pokedex[pokemonName]

	if !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Stats: \n")

		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}

		fmt.Printf("Types: \n")

		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
	}

	return nil
}

func commandPokedex(c *config) error {
	fmt.Println("Your Pokedex:")

	for name := range pokedex {
		fmt.Printf("- %s\n", name)
	}

	return nil
}

func getSupportedCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays a map of the region",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map back",
			description: "Displays a map of the region",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of pokemon in a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays your caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c := &config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: "",
		Params:   []string{},
	}

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		text := scanner.Text()
		cleanedInput := cleanInput(text)

		if len(cleanedInput) == 0 {
			continue
		}

		c.Params = cleanedInput[1:]

		commandName := cleanedInput[0]
		commandMap := getSupportedCommands()
		command, ok := commandMap[commandName]

		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(c)
		if err != nil {
			fmt.Println(err)
		}
	}
}
