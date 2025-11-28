package repl

import (
	"fmt"
	"github.com/CSelvidge/pokedexcli/internal/pokeapi"
	"github.com/CSelvidge/pokedexcli/internal/actors"
	"os"
	"strings"
	"math/rand"
	"time"
	"encoding/json"
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
	cfg.currentLocation = locationName
	cfg.currentLocationURL = url
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

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		fmt.Println("Please provide a Pokemon name to catch. Usage: catch <pokemon-name>")
		return nil
	}

	
	pokemonName := args[0]
	exists, err := pokemonStreamingCheck(cfg, pokemonName)
	if !exists {
		return fmt.Errorf("Pokemon is not in this area")
	}
	if err != nil {
		fmt.Printf("Error checking byte stream for pokemon\n")
		return err
	}

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	
	pokemon:= &actors.Pokemon{}
	if err := pokeapi.GenericURLCaller(url, cfg.cache, pokemon); err != nil {
		fmt.Printf("Error fetching Pokemon data: %v\n", err)
		return err
	}

	rand.Seed(time.Now().UnixNano())
	catchChance := rand.Intn(300)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if catchChance > pokemon.BaseExperience {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	} else {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.user.CaughtPokemon[pokemon.Name] = *pokemon
	}
	return nil
}

func pokemonStreamingCheck(cfg *config, args ...string) (bool, error) {
	if len(args) == 0 {
		return false, fmt.Errorf("No pokemon name provided")
	}


	cacheData, exists := cfg.cache.Get(cfg.currentLocationURL)
	if !exists {
		return false, nil
	}

	var partialResponse struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			}`json:"pokemon"`
		}`json:"pokemon_encounters"`
	}

	err := json.Unmarshal(cacheData, &partialResponse)
	if err != nil {
		return false, fmt.Errorf("Failed to parse cache data: %w", err)
	}

	for _, encounter := range partialResponse.PokemonEncounters {
		if strings.EqualFold(encounter.Pokemon.Name, args[0]) {
			return true, nil
		}
	
	}
	return false, nil
}

func commandInspect(cfg *config, args ...string) error {
	pokemonName := args[0]

	pokemon, exists := cfg.user.CaughtPokemon[pokemonName]
	if !exists {
		return  fmt.Errorf("You have not caught a %s\n", pokemonName)
	}

	fmt.Printf("Name: %s\nAbilities:\n", pokemon.Name)
	for _, ability := range pokemon.Abilities {
		fmt.Printf(" - %s", ability.Ability.Name)
	}
	fmt.Printf("\nMoves:\n")
	for _, move := range pokemon.Moves {
		fmt.Printf(" - %s\n", move.Move.Name)
	}
	return nil
}

func commandFullInspect(cfg *config, args ...string) error {
	if err := commandInspect(cfg, args[0]); err != nil {
		return err
	}

	pokemon := cfg.user.CaughtPokemon[args[0]]

	fmt.Printf("\nTypes:\n")
	for _, kind := range pokemon.Types {
		fmt.Printf(" - %s", kind.Type.Name)
	}
	fmt.Printf("\nStats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s", stat.Stat.Name)
	}
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.user.CaughtPokemon) <= 0 {
		return fmt.Errorf("You have not caught any pokemon.")
	}
	fmt.Printf("Your pokedex:\n")
	for _, pokemon := range cfg.user.CaughtPokemon {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
