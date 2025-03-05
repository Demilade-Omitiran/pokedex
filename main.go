package main

import (
	"fmt"

	"bufio"

	"os"

	"strings"

	"internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	c := &config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		// scanner.Scan()
		text := scanner.Text()
		cleanedInput := cleanInput(text)

		if len(cleanedInput) == 0 {
			continue
		}

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
