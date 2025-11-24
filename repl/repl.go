package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*locationResponse) error
}

type locationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var commandDictionary = make(map[string]cliCommand)
var config = &locationResponse{}

func Start() {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Type 'help' to see available commands.")
	initMap()
	getUserInput()
}

func initMap() {
	commandDictionary["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex CLI",
		callback:    commandExit,
	}
	commandDictionary["help"] = cliCommand{
		name:        "help",
		description: "List all available commands",
		callback:    commandHelp,
	}
	commandDictionary["map"] = cliCommand{
		name:        "map",
		description: "Show locations to travel to",
		callback:    commandMap,
	}
	commandDictionary["mapb"] = cliCommand{
		name:        "mapb",
		description: "Show previous locations",
		callback:    commandMapb,
	}
}

func getUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Pokedex >")
	for scanner.Scan() {
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			fmt.Println("Please enter at least one character")
			fmt.Printf("Pokedex >")
			continue
		}
		commandName := input[0]
		command, exists := commandDictionary[commandName]
		if !exists {
			fmt.Printf("Unknown command: %s\n", commandName)
			fmt.Printf("Pokedex >")
			continue
		}
		command.callback(config)
		fmt.Printf("Pokedex >")
	}

}

func cleanInput(input string) []string {
	lowered := strings.ToLower(input)
	trimmed := strings.TrimSpace(lowered)
	words := strings.Fields(trimmed)

	return words
}

func commandExit(cfg *locationResponse) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *locationResponse) error {
	fmt.Println("Available commands:")
	for _, cmd := range commandDictionary {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *locationResponse) error {
	url := ""
	if cfg.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = cfg.Next
	}

	if err := genericURLCaller(url, cfg); err != nil {
		fmt.Printf("Error fetching locations: %v\n", err)
		return err
	}

	for _, location := range cfg.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandMapb(cfg *locationResponse) error {
	url := ""
	if cfg.Previous == "" {
		fmt.Println("No previous locations available. You must advance at least once first.")
		return nil
	} else {
		url = cfg.Previous
	}

	if err := genericURLCaller(url, cfg); err != nil {
		fmt.Printf("Error fetching locations: %v\n", err)
		return err
	}

	for _, location := range cfg.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func genericURLCaller(url string, cfg *locationResponse) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(cfg); err != nil {
		return err
	}

	return nil
}
