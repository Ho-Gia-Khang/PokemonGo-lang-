package main

import (
	"os"

	"github.com/go-jose/go-jose/v3/json"
)

func main() {
	crawlPokemonsDriver(numberOfPokemons)
	data, _ := os.ReadFile("pokedex.json")
	allPokemons := []Pokemon{}
	json.Unmarshal(data, &allPokemons)

}
