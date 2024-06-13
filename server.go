package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
type Pokedex struct {
	Pokemon     []Pokemon `json:"Pokemon"`
	CoordinateX int
	CoordinateY int
}

type Player struct {
	Name              string    `json:"Name"`
	ID                string    `json:"ID"`
	PlayerCoordinateX int       `json:"PlayerCoordinateX"`
	PlayerCoordinateY int       `json:"PlayerCoordinateY"`
	Inventory         []Pokemon `json:"Inventory"`
	IsTurn            bool
	Addr              *net.UDPAddr
	sync.Mutex
}

var players = make(map[string]*Player)
var pokeDexWorld = make(map[string]*Pokedex)
var Pokeworld [20][20]string
var inventory1 = make(map[string]*Pokemon)
var inventory []Pokemon

func randomInt(max int64) (int64, error) {
	// Generate a random big integer in the range [0, max)
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err // Return the error if any
	}
	return n.Int64(), nil // Convert the big integer to int64 and return
}
func PassingPokemontoInventory(pokemon *Pokemon, player *Player) {
	player.Lock() // Lock the player instance
	defer player.Unlock()
	player.Inventory = append(player.Inventory, *pokemon)
}
func PassingPlayertoJson(filename string, player *Player) {

	updatedData, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
	}
	if err := os.WriteFile(filename, updatedData, 0666); err != nil {
		fmt.Println("Error:", err)
	}

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

	max := int64(19) // Example maximum value
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

func CheckForPokemonEncounter(player *Player, pokemon *Pokedex) {
	for _, pokedex := range pokemon.Pokemon {
		if player.PlayerCoordinateX == pokemon.CoordinateX && player.PlayerCoordinateY == pokemon.CoordinateY {
			PassingPokemontoInventory(&pokedex, player)
			fmt.Println("Pokemon encountered:", pokedex.Name)

		}
	}
}

func printWorld(x, y int) string {
	world := "" // Initialize the world as an empty string
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			// If the current position matches the player's coordinates
			if i == x && j == y {
				world += "P"

			} else if Pokeworld[i][j] == "M" {
				world += "M" // Append "M" for Pokemon
			} else {
				world += "-" // Append "-" for Empty space
			}
		}
		world += "\n" // New line after each row
	}
	return world
}

func PokeCat(Id string, playername string, x int, y int, conn *net.UDPConn, Addr *net.UDPAddr) string {
	// Check if the coordinates are within the bounds of Pokeworld.
	if x >= 0 && x < 20 && y >= 0 && y < 20 {
		// Check if the position is already occupied.
		if Pokeworld[x][y] == "" || Pokeworld[x][y] == "M" {
			// Place the player at the specified coordinates.
			Pokeworld[x][y] = Id
			if player, exists := players[Id]; exists {
				// Player exists, update the existing player's fields
				player.Name = playername
				player.PlayerCoordinateX = x
				player.PlayerCoordinateY = y
				player.Addr = Addr
				for _, pokedex := range pokeDexWorld {
					CheckForPokemonEncounter(players[Id], pokedex)
				}
				fileName := fmt.Sprintf("pokeInventory%s.json", Id)
				PassingPlayertoJson(fileName, players[Id])
			} else {
				// Player does not exist, create a new one
				players[Id] = &Player{
					Name:              playername,
					ID:                Id,
					PlayerCoordinateX: x,
					PlayerCoordinateY: y,
					Addr:              Addr,
				}
				for _, pokedex := range pokeDexWorld {
					CheckForPokemonEncounter(players[Id], pokedex)
				}
				fileName := fmt.Sprintf("pokeInventory%s.json", Id)
				PassingPlayertoJson(fileName, players[Id])
			}

			fmt.Println("Player placed at", x, y)
			world := printWorld(x, y)
			return world
		} else {
			battleScene(players[Id], players[Pokeworld[x][y]], conn, Addr, players[Pokeworld[x][y]].Addr)
			return "Battle"
		}
	}
	return ""
}

func movePlayer(idStr string, direction string, conn *net.UDPConn) string {
	player, exists := players[idStr]
	if !exists {
		fmt.Println("Player does not exist.")

	}
	deltaX := map[string]int{"Up": -1, "Down": 1}[direction]
	newX := player.PlayerCoordinateX + deltaX
	deltaY := map[string]int{"Left": -1, "Right": 1}[direction]
	newY := player.PlayerCoordinateY + deltaY
	Pokeworld[player.PlayerCoordinateX][player.PlayerCoordinateY] = ""

	PokeK := PokeCat(idStr, player.Name, newX, newY, conn, player.Addr)
	return PokeK

}

func main() {

	for k := 0; k < 20; k++ {
		pokemon, err := getRandomPokemon("pokedex.json")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		pokedex := Pokedex{Pokemon: []Pokemon{*pokemon}}
		key := strconv.Itoa(k)
		pokeDexWorld[key] = &pokedex
		positionofPok(&pokedex)
		Pokeworld[pokedex.CoordinateX][pokedex.CoordinateY] = "M"
		fmt.Println("Pokemon:", pokemon.Name, "X:", pokedex.CoordinateX, "Y:", pokedex.CoordinateY)
	}
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
			players[idStr] = &Player{Name: parts[1], ID: idStr}
			xBigInt, _ := rand.Int(rand.Reader, big.NewInt(10))
			yBigInt, _ := rand.Int(rand.Reader, big.NewInt(10))
			x := int(xBigInt.Int64())
			y := int(yBigInt.Int64())

			PokeC := PokeCat(idStr, parts[1], x, y, conn, addr)
			connectclient := fmt.Sprintf("Client connected: %s %s ID: %s", parts[1], addr, idStr)
			_, err := conn.WriteToUDP([]byte(connectclient), addr)
			if err != nil {
				fmt.Println("Error sending connect message to client:", err)
			}
			_, err = conn.WriteToUDP([]byte(PokeC), addr)
			if err != nil {
				fmt.Println("Error sending connect message to client:", err)
			}
			// Handle connection...
		case "Info":
			Info := fmt.Sprintln("Player Info:%s", idStr)
			_, err := conn.WriteToUDP([]byte(Info), addr)
			if err != nil {
				fmt.Println("Error sending connect message to client:", err)
			}
			// Display player info...
		case "DISCONNECT":
			fmt.Println("Disconnected from server.")
			return
		case "Up", "Down", "Left", "Right":
			PokeK := movePlayer(idStr, parts[0], conn)
			fmt.Println(PokeK)
			_, err := conn.WriteToUDP([]byte(PokeK), addr)
			if err != nil {
				fmt.Println("Error sending connect message to client:", err)
			}
		case "Inventory":
			for _, inv := range players[idStr].Inventory {
				inventoryDetails := fmt.Sprintf("Player Inventory: Name: %s, Level: %d", inv.Name, inv.Level)
				_, err := conn.WriteToUDP([]byte(inventoryDetails), addr)
				if err != nil {
					fmt.Println("Error sending connect message to client:", err)
				}
			}
		}

	}
}
