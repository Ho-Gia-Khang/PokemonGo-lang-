package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
)

type Stats struct {
	HP         float32 `json:"HP"`
	Attack     float32 `json:"Attack"`
	Defense    float32 `json:"Defense"`
	Speed      int     `json:"Speed"`
	Sp_Attack  int     `json:"Sp_Attack"`
	Sp_Defense int     `json:"Sp_Defense"`
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
}

type Pokedex struct {
	MyPokemon   []Pokemon `json:"Pokemon"`
	CoordinateX int
	CoordinateY int
}

type Player struct {
	Name              string
	ID                string
	PlayerCoordinateX int
	PlayerCoordinateY int
	Addr              *net.UDPAddr
}

var players = make(map[string]*Player)
var Pokeworld [200][200]string

func randomInt(max int64) (int64, error) {
	// Generate a random big integer in the range [0, max)
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err // Return the error if any
	}
	return n.Int64(), nil // Convert the big integer to int64 and return
}
func getRandomPokemon(filename string) (*Pokemon, error) {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of Pokemon structs
	var pokemons []Pokemon
	err = json.Unmarshal(data, &pokemons)
	if err != nil {
		return nil, err
	}

	// Generate a random index
	index, err := randomInt(int64(len(pokemons)))
	if err != nil {
		return nil, err
	}

	// Return the randomly selected Pokemon
	return &pokemons[index], nil

}

func positionofPok(pokedex *Pokedex) {
	max := int64(199) // Example maximum value
	x, err := randomInt(max)
	if err != nil {
		fmt.Println("Error generating random x:", err)

	}
	y, err := randomInt(max)
	if err != nil {
		fmt.Println("Error generating random y:", err)

	}
	pokedex.CoordinateX = int(x)
	pokedex.CoordinateY = int(y)
}

func printWorld(x, y int) {
	for k := 0; k < 50; k++ {
		pokemon, err := getRandomPokemon("pokedex.json")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		pokedex := Pokedex{MyPokemon: []Pokemon{*pokemon}}
		positionofPok(&pokedex)
		Pokeworld[pokedex.CoordinateX][pokedex.CoordinateY] = pokemon.Name
		fmt.Println(pokedex.MyPokemon[0].Name, "placed at", pokedex.CoordinateX, pokedex.CoordinateY)

	}

	for i := 0; i < len(Pokeworld); i++ {
		for j := 0; j < len(Pokeworld[i]); j++ {
			// If the current position matches the player's coordinates
			if i == y && j == x {
				fmt.Print("P") // Print "P" for Player
			} else {
				// Print the value at the current position
				// Here, you might customize the printing based on what each value represents
				switch Pokeworld[i][j] {
				case "":
					fmt.Print("0") // Empty space
				default:
					fmt.Print("M") // Unknown
				}
			}
		}
		fmt.Println() // New line after each row
	}
}

func PokeCat(Id string, playername string, x int, y int) {
	// Check if the coordinates are within the bounds of Pokeworld.
	if x >= 0 && x < 200 && y >= 0 && y < 200 {
		// Check if the position is already occupied.
		if Pokeworld[x][y] == "" {
			// Place the player at the specified coordinates.
			Pokeworld[x][y] = Id
			players[Id] = &Player{
				Name:              playername,
				ID:                Id,
				PlayerCoordinateX: x,
				PlayerCoordinateY: y,
				Addr:              players[Id].Addr,
			}
			printWorld(x, y)
			fmt.Println("Player placed at", x, y)
		} else {
			// Handle the case where the position is already occupied.
			fmt.Println("Position is already occupied.")
		}
	} else {
		// Handle the case where the coordinates are out of bounds.s
	}
}
func PokeBat() {

}
func movePlayer(idStr string, direction string) {
	player, exists := players[idStr]
	if !exists {
		fmt.Println("Player does not exist.")
		return
	}
	deltaY := map[string]int{"Up": -1, "Down": 1}[direction]
	newY := player.PlayerCoordinateY + deltaY
	deltaX := map[string]int{"Left": -1, "Right": 1}[direction]
	newX := player.PlayerCoordinateX + deltaX
	Pokeworld[player.PlayerCoordinateX][player.PlayerCoordinateY] = ""
	PokeCat(idStr, player.Name, newX, newY)
}

func main() {

	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			panic(err)
		}

		clientAddr := (strings.Replace(addr.String(), ".", "", -1)) // This removes all periods      // This includes IP and port
		clientAddr = (strings.Replace(clientAddr, ":", "", -1))     // This removes the colon         // This includes IP and port
		clientAddr = (strings.Replace(clientAddr, " ", "", -1))     // This removes all spaces        // This includes IP and port                // This converts the length to a string // This includes IP and port

		idStr := clientAddr

		commands := string(buffer[:n])
		parts := strings.Split(commands, " ")

		switch parts[0] {
		case "CONNECT":
			fmt.Println("Unique ID Int:", idStr)
			players[idStr] = &Player{Name: parts[1], Addr: addr, ID: idStr}
			xBigInt, _ := rand.Int(rand.Reader, big.NewInt(200))
			yBigInt, _ := rand.Int(rand.Reader, big.NewInt(200))
			x := int(xBigInt.Int64())
			y := int(yBigInt.Int64())
			PokeCat(idStr, parts[1], x, y)
			fmt.Println("Client connected:", parts[1], addr, "ID:", idStr)
			fmt.Println("Player placed at", x, y)
			// Handle connection...
		case "Info":
			fmt.Println("Player Info:", idStr)
			// Display player info...
		case "DISCONNECT":
			fmt.Println("Disconnected from server.")
			return
		case "Up", "Down", "Left", "Right":
			movePlayer(idStr, parts[0])
		}

	}
}
