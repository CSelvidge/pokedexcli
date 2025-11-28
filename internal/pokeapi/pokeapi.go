package pokeapi

import (
	"encoding/json"
	"github.com/CSelvidge/pokedexcli/internal/pokecache"
	"net/http"
)

func GenericURLCaller(url string, cache *pokecache.Cache, target interface{}) error { //generic function to fill different types of structs

	val, exists := cache.Get(url) // check that cache first!
	if exists {
		if err := json.Unmarshal(val, target); err != nil {
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
	if err := dec.Decode(target); err != nil {
		return err
	}

	data, err := json.Marshal(target)
	if err != nil {
		return err
	}

	cache.Add(url, data)
	return nil
}
