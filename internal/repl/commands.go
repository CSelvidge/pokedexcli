package repl

import (
	"fmt"
	"github.com/CSelvidge/pokedexcli/internal/pokeapi"
	"os"
)

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Available commands:")
	for _, cmd := range commandDictionary {
		fmt.Printf(" - %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	var locationMap locationResponse

	url := ""
	if cfg.nextLocationsURL == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = cfg.nextLocationsURL
	}

	if err := pokeapi.GenericURLCaller(url, cfg.cache, &locationMap); err != nil {
		fmt.Printf("Error fetching locations: %v\n", err)
		return err
	}

	cfg.nextLocationsURL = locationMap.Next
	cfg.previousLocationsURL = locationMap.Previous

	for _, location := range locationMap.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	var locationMap locationResponse

	url := ""
	if cfg.previousLocationsURL == "" {
		fmt.Println("No previous locations available. You must advance at least once first.")
		return nil
	} else {
		url = cfg.previousLocationsURL
	}

	if err := pokeapi.GenericURLCaller(url, cfg.cache, &locationMap); err != nil {
		fmt.Printf("Error fetching locations: %v\n", err)
		return err
	}

	cfg.previousLocationsURL = locationMap.Previous
	cfg.nextLocationsURL = locationMap.Next

	for _, location := range locationMap.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	var locationInfo exploreResponse
	foundPokemon := []string{}
	if len(args) == 0 {
		fmt.Println("Please provide a location name to explore. Usage: explore <location-name>")
		return nil
	}
	locationName := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	if err := pokeapi.GenericURLCaller(url, cfg.cache, &locationInfo); err != nil {
		fmt.Printf("Error exploring location: %v\n", err)
		return err
	}
	fmt.Printf("Exploring %s...\n", locationInfo.Name)
	for _, encounter := range locationInfo.PokemonEncounters {
		foundPokemon = append(foundPokemon, encounter.Pokemon.Name)
	}
	if len(foundPokemon) == 0 {
		fmt.Printf("No Pokemon found in %s.\n", locationInfo.Name)
		return nil
	}
	fmt.Printf("Found Pokemon:\n")
	for _, pokemon := range foundPokemon {
		fmt.Printf(" - %s\n", pokemon)
	}
	return nil

}
