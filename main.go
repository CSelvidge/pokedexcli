package main

import (
	"fmt"
	"github.com/CSelvidge/pokedexcli/internal/pokecache"
	"github.com/CSelvidge/pokedexcli/internal/repl"
	"os"
)

func main() {
	cache, err := initCache()
	if err != nil {
		fmt.Printf("Error initializing cache: %v\n", err)
		os.Exit(1)
	}
	repl.Start(cache)
}

func initCache() (*pokecache.Cache, error) {
	var err error
	cache, err := pokecache.NewCache(repl.GetCacheSettings())
	if err != nil {
		fmt.Printf("Error initializing cache: %v\n", err)
		return nil, err
	}
	return cache, nil
}
