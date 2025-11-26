package repl

import (
	"bufio"
	"github.com/CSelvidge/pokedexcli/internal/pokecache"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"strconv"
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
var cache = &pokecache.Cache{}

func Start() {
	fmt.Println("Welcome to the Pokedex!")
	initCache()
	fmt.Println("Type 'help' to see available commands.")
	initMap()
	getUserInput()
}

func initCache() {
var err error
cache, err = pokecache.NewCache(getCacheSettings())
if err != nil {
	fmt.Printf("Error initializing cache: %v\n", err)
os.Exit(1)
}
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


func getCacheSettings() (Type string, Life int) {
	durationType := ""
	durationLife := 0
	plural := "s"

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter cache duration type (eg: Seconds, Minutes, Hours) below:\n")
	fmt.Printf("pokedex >")
	for scanner.Scan() {
		durationInput := cleanInput(scanner.Text())
		if len(durationInput) == 0 {
			fmt.Printf("\nPlease enter at least one character for duration type, the first character is accepted.\n>")
			continue
		}

		switch durationInput[0] {
		case "seconds", "second", "s":
			fmt.Printf("Duration type set to Seconds\n")
			durationType = "second"
		case "minutes", "minute", "m":
			fmt.Printf("Duration type set to Minutes\n")
			durationType = "minute"
		case "hours", "hour", "h":
			fmt.Printf("Duration type set to Hours, however this type is for fun and cache should not live this long.\n")
			durationType = "hour"
		default:
			fmt.Printf("Invalid duration type. Please enter Seconds, Minutes, or Hours.\n>")
			continue
		}
		break
	}

	fmt.Printf("Enter cache life duration in number of %s below, only integers are accepted:\n", durationType)
	fmt.Printf("pokedex >")
	for scanner.Scan() {
		lifeInput := cleanInput(scanner.Text())
		if len(lifeInput) == 0 {
			fmt.Println("Please enter a valid duration number.")
			continue
		}
		firstInt := lifeInput[0]
		num, err := strconv.Atoi(firstInt)
		if err != nil || num <= 0 {
			fmt.Println("Invalid duration number. Please enter a positive integer of at least 1.")
			continue
		}
		durationLife = num
		if num == 1 {
			plural = ""
		}
		fmt.Printf("%d %s%s set for cache duration life.\n", durationLife, durationType, plural)
	break
	}
	return durationType, durationLife
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
	val, exists := cache.Get(url)
	if exists {
		if err := json.Unmarshal(val, cfg); err != nil {
			return err
		}
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(cfg); err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	cache.Add(url, data)
	return nil
}
