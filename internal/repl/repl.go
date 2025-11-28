package repl

import (
	"bufio"
	"fmt"
	"github.com/CSelvidge/pokedexcli/internal/pokecache"
	"github.com/CSelvidge/pokedexcli/internal/actors"
	"os"
	"strconv"
	"strings"
)

var commandDictionary = make(map[string]cliCommand)

func Start(cache *pokecache.Cache, user *actors.User) {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Type 'help' to see available commands.")
	initMap()
	cfg := newConfig(cache, user)
	getUserInput(cfg)
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
	commandDictionary["explore"] = cliCommand{
		name:        "explore",
		description: "Explore location to find Pokemon! Usage is `explore <location-name>`",
		callback:    commandExplore,
	}
	commandDictionary["catch"] = cliCommand{
		name:        "catch",
		description: "Catch a Pokemon! Usage is `catch <pokemon-name>`",
		callback:    commandCatch,
	}
	commandDictionary["inspect"] = cliCommand{
		name: "inspect",
		description: "Brief inspection of caught pokemon Usage is `inspect <pokemon-name>`",
		callback: commandInspect,
	}
	commandDictionary["fullinspect"] = cliCommand{
		name: "fullinspect",
		description: "Inspect all stored stats for caught pokemon. Usage is same as inspect.",
		callback: commandFullInspect,
	}
	commandDictionary["pokedex"] = cliCommand{
		name: "pokedex",
		description: "list all caught pokemon",
		callback: commandPokedex,
	}
}

func newConfig(cache *pokecache.Cache, user *actors.User) *config {
	cfg := &config{
		nextLocationsURL:     "",
		previousLocationsURL: "",
		currentLocation:      "",
		currentLocationURL:   "",
		cache:                cache,
		user:                 user,
	}
	return cfg
}

func GetCacheSettings() (Type string, Life int) {
	durationType := ""
	durationLife := 0
	plural := "s"

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter cache duration type (eg: Seconds, Minutes, Hours) below:\n")
	fmt.Printf("Pokedex >")
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
	fmt.Printf("Pokedex >")
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

func getUserInput(cfg *config) {
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
		commandArgs := input[1:] //Future proofing for multiple arguments

		command, exists := commandDictionary[commandName]

		if !exists {
			fmt.Printf("Unknown command: %s\n", commandName)
			fmt.Printf("Pokedex >")
			continue
		}
		executeCommand(cfg, command, commandArgs)
		fmt.Printf("Pokedex >")
	}
}

func cleanInput(input string) []string {
	lowered := strings.ToLower(input)
	trimmed := strings.TrimSpace(lowered)
	words := strings.Fields(trimmed)

	return words
}

func executeCommand(cfg *config, command cliCommand, args []string) { //input sanitized in function that called, so we know command is valid
	var arguments string //if multiple aruguments are needed in the future this will need to be changed to a slice
	if len(args) > 0 {
		arguments = args[0] //should only ever get 1 argument, but this allows for easier future proofing if more args are needed
	}

	err := command.callback(cfg, arguments) //functions are variadic, so arguments can be empty
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
