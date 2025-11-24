package repl

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"encoding/json"
)

type cliCommand struct {
	name        string
	description string
	config	  	*locationResponse
	callback    func() error
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
	initConfiguration()
	initMap()
	getUserInput()
}

func initConfiguration() {
	config.Next = ""
	config.Previous = ""
	//Unset until needed
}

func initMap() {
	commandDictionary["exit"] = cliCommand{
		name: "exit",
		description: "Exit the Pokedex CLI",
		config: config,
		callback: commandExit,
	}
	commandDictionary["help"] = cliCommand{
		name: "help",
		description: "List all available commands",
		config: config,
		callback: commandHelp,
	}
	commandDictionary["map"] = cliCommand{
		name: "map",
		description: "Show locations to travel to",
		config: config,
		callback: commandMap,
	}
	commandDictionary["mapb"] = cliCommand{
		name: "mapb",
		description: "Show previous locations",
		config: config,
		callback: commandMapb,
	}
}

func getUserInput() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Pokedex >")
		for scanner.Scan(){		
			input := cleanInput(scanner.Text())
			if len(input) == 0 {
				fmt.Println("Please enter atleast one character")
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
			command.callback()
			fmt.Printf("Pokedex >")
		}
	
}

func cleanInput(input string) []string {
	lowered := strings.ToLower(input)
	trimmed := strings.TrimSpace(lowered)
	words := strings.Fields(trimmed)

    return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Available commands:")
	for _, cmd := range commandDictionary {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap() error {
	url := ""
	if config.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = config.Next
	}
	
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&config); err != nil {
		fmt.Println("Error decoding map data", err)
		return err
	}

	for _, location := range config.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandMapb() error {
	url := ""
	if config.Previous == "" {
		fmt.Println("No previous locations available. You must advance at least once first.")
		return nil
	}else {
		url = config.Previous
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching map data:", err)
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&config); err != nil {
		fmt.Println("Error decoding map data", err)
		return err
	}

	for _, location := range config.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}
