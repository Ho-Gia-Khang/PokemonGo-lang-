package main

// import (
// 	"fmt"
// 	"net/http"
// )

type DamegeWhenAttacked struct {
	Element     string  
	Coefficient float64 
}

type Stats struct {
	HP                 int                  
	Defense            int                  
	Speed              int                  
	Sp_Attack          int                  
	Sp_Defense         int                  
}

type GenderRatio struct {
	MaleRatio  int 
	FemaleRatio int 
}

type Profile struct {
	Name 			string               
	Weight          float64              
	CatchRate 	 int                  
	GenderRatio     GenderRatio 
	EggGroup        []string 
	HatchSteps	  int
	Abilities	   []string 
}

type NaturalMoves struct {
	Name string
	Element string
	Power int
	Acc 	int
	PP 		int
	Description string
}

type MachineMoves struct {
	Name string
	Element string
	Power int
	Acc 	int
	PP 		int
	Description string
}

type TutorMoves struct {
	Name string
	Element string
	Power int
	Acc 	int
	PP 		int
	Description string
}

type EggMoves struct {
	Name string
	Element string
	Power int
	Acc 	int
	PP 		int
	Description string
}

type Moves struct {
	NaturalMoves []NaturalMoves
	MachineMoves []MachineMoves
	TutorMoves   []TutorMoves
	EggMoves     []EggMoves
}

type Pokemon struct {
	Elements           []string             
	EV                 int                                
	Profile			Profile              
	DamegeWhenAttacked []DamegeWhenAttacked 
	EvolutionLevel    int
	Moves			  Moves                  
}

// func fetchPokemons(index int) []Pokemon {
// 	url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", index);
// 		resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println("Error fetching genre page:", err)
// 		return nil
// 	}
// 	defer resp.Body.Close()
// }