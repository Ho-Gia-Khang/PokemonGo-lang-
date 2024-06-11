package main

import (
	"os"

	"github.com/go-jose/go-jose/v3/json"
)

type Stats struct {
	HP         float32 `json:"HP"`
	Attack     float32 `json:"Attack"`
	Defense    float32 `json:"Defense"`
	Speed      int     `json:"Speed"`
	Sp_Attack  float32 `json:"Sp_Attack"`
	Sp_Defense float32 `json:"Sp_Defense"`
}

type GenderRatio struct {
	MaleRatio   float32 `json:"MaleRatio"`
	FemaleRatio float32 `json:"FemaleRatio"`
}

type Profile struct {
	Height      float32     `json:"Height"`
	Weight      float32     `json:"Weight"`
	CatchRate   float32     `json:"CatchRate"`
	GenderRatio GenderRatio `json:"GenderRatio"`
	EggGroup    string      `json:"EggGroup"`
	HatchSteps  int         `json:"HatchSteps"`
	Abilities   string      `json:"Abilities"`
}

type DamegeWhenAttacked struct {
	Element     string  `json:"Element"`
	Coefficient float32 `json:"Coefficient"`
}

type Moves struct {
	Name    string   `json:"Name"`
	Element []string `json:"Element"`
	Power   float32  `json:"Power"`
	Acc     int      `json:"Acc"`
	PP      int      `json:"PP"`
}

type Pokemon struct {
	Name               string               `json:"Name"`
	Elements           []string             `json:"Elements"`
	EV                 int                  `json:"EV"`
	Stats              Stats                `json:"Stats"`
	Profile            Profile              `json:"Profile"`
	DamegeWhenAttacked []DamegeWhenAttacked `json:"DamegeWhenAttacked"`
	EvolutionLevel     int                  `json:"EvolutionLevel"`
	NextEvolution      string               `json:"NextEvolution"`
	Moves              []Moves              `json:"Moves"`
	Experience         int                  `json:"Experience"`
	Level              int                  `json:"Level"`
}

func main() {
	// crawlPokemonsDriver(numberOfPokemons)
	data, _ := os.ReadFile("pokedex.json")
	allPokemons := []Pokemon{}
	json.Unmarshal(data, &allPokemons)

	// player 1's pokemons
	Venusaur := allPokemons[2]
	Charizard := allPokemons[5]
	Blastoise := allPokemons[8]
	Pidgeot := allPokemons[17]
	Rattata := allPokemons[18]
	player1 := Player{
		Inventory: []Pokemon{Venusaur, Charizard, Blastoise, Pidgeot, Rattata},
	}

	// player 2's pokemons
	Butterfree := allPokemons[11]
	Pikachu := allPokemons[24]
	Sandslash := allPokemons[27]
	Nidoqueen := allPokemons[30]
	Ninetales := allPokemons[37]
	player2 := Player{
		Inventory: []Pokemon{Butterfree, Pikachu, Sandslash, Nidoqueen, Ninetales},
	}

	battleScene(&player1, &player2)
}
