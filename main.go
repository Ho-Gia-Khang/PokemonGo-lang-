// package main

// import (
// 	"os"

// 	"github.com/go-jose/go-jose/v3/json"
// )

// func main() {
// 	// crawlPokemonsDriver(numberOfPokemons)
// 	data, _ := os.ReadFile("pokedex.json")
// 	allPokemons := []Pokemon{}
// 	json.Unmarshal(data, &allPokemons)

// 	// player 1's pokemons
// 	Venusaur := allPokemons[2]
// 	Charizard := allPokemons[5]
// 	Blastoise := allPokemons[8]
// 	Pidgeot := allPokemons[17]
// 	Rattata := allPokemons[18]
// 	player1 := Player{
// 		Inventory: []Pokemon{Venusaur, Charizard, Blastoise, Pidgeot, Rattata},
// 	}

// 	// player 2's pokemons
// 	Butterfree := allPokemons[11]
// 	Pikachu := allPokemons[24]
// 	Sandslash := allPokemons[27]
// 	Nidoqueen := allPokemons[30]
// 	Ninetales := allPokemons[37]
// 	player2 := Player{
// 		Inventory: []Pokemon{Butterfree, Pikachu, Sandslash, Nidoqueen, Ninetales},
// 	}

// 	battleScene(&player1, &player2)
// }
