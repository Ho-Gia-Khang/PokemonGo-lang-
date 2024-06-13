package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
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
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var players []Player
	if err := json.Unmarshal(data, &players); err != nil {
		fmt.Println("Error:", err)
	}

	players = append(players, *player)
	updatedData, err := json.MarshalIndent(players, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
	}
	if err := os.WriteFile(filename, updatedData, 0644); err != nil {
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
				PassingPlayertoJson("pokeInventory.json", players[Id])
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
				PassingPlayertoJson("pokeInventory.json", players[Id])
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

// var (
// 	mutex          sync.Mutex
// 	cond           = sync.NewCond(&mutex)
// 	playernum      int
// 	numberofPlayer []*Player
// )

// func PokeBat(idStr string, conn *net.UDPConn) {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	player, exists := players[idStr]
// 	if !exists {
// 		fmt.Println("Player does not exist.")
// 		return
// 	}
// 	playernum += 1
// 	numberofPlayer = append(numberofPlayer, player)

// 	if playernum < 2 {
// 		// Wait with a timeout to avoid indefinite waiting
// 		waitTimeout(cond, &mutex, time.Second*20)
// 	} else if playernum == 2 {
// 		cond.Broadcast()
// 		go battleScene(numberofPlayer[0], numberofPlayer[1], conn, numberofPlayer[0].Addr, numberofPlayer[1].Addr)
// 	}
// }

// func waitTimeout(cond *sync.Cond, mutex *sync.Mutex, timeout time.Duration) {
// 	mutex.Lock()
// 	defer mutex.Unlock()

//		var timeoutCh = time.After(timeout)
//		go func() {
//			cond.Wait()
//		}()
//		select {
//		case <-timeoutCh:
//			// Timeout occurred
//			fmt.Println("Timeout waiting for another player.")
//			return
//		default:
//			// Condition was signaled before timeout
//		}
//	}
func battleScene(player1 *Player, player2 *Player, conn *net.UDPConn, addr1, addr2 *net.UDPAddr) {

	if player1 == nil {
		fmt.Println("Error: player1 is nil")
		return
	}
	if conn == nil {
		fmt.Println("Error: conn is nil")
		return
	}
	if addr1 == nil {
		fmt.Println("Error: addr1 is nil")
		return
	}
	if len(player1.Inventory) < 3 {
		fmt.Println("Player 1 has less than 3 pokemons")
		conn.WriteToUDP([]byte("You have less than 3 pokemons"), addr1)
		return
	} else if len(player2.Inventory) < 3 {
		fmt.Println("Player 2 has less than 3 pokemons")
		conn.WriteToUDP([]byte("You have less than 3 pokemons"), addr2)
		return
	}

	// Player 1 select 3 Pokemons
	fmt.Println("Player 1 please select 3 pokemons from:")
	conn.WriteToUDP([]byte("Player 1 please select 3 pokemons from:\n"), addr1)
	for i := range player1.Inventory {
		printPokemonInfo(i, player1.Inventory[i])
		conn.WriteToUDP([]byte(fmt.Sprintf("%d: %s\n", i, player1.Inventory[i].Name)), addr1)
	}
	player1Pokemons := selectPokemon(player1, conn, addr1)

	// Player 2 select 3 Pokemons
	fmt.Println("Player 2 please select 3 pokemons from:")
	conn.WriteToUDP([]byte("Player 2 please select 3 pokemons from:\n"), addr2)
	for i := range player2.Inventory {
		printPokemonInfo(i, player2.Inventory[i])
		conn.WriteToUDP([]byte(fmt.Sprintf("%d: %s\n", i, player2.Inventory[i].Name)), addr2)
	}
	player2Pokemons := selectPokemon(player2, conn, addr2)

	allBattlingPokemons := append(*player1Pokemons, *player2Pokemons...)
	firstAttacker := getFirstAttacker(allBattlingPokemons)
	var firstDefender *Pokemon

	fmt.Println("Battle start!")
	conn.WriteToUDP([]byte("Battle start!\n"), addr1)
	conn.WriteToUDP([]byte("Battle start!\n"), addr2)

	if isContain(*player1Pokemons, *firstAttacker) {
		firstDefender = getFirstDefender(*player2Pokemons)
		fmt.Println("Player 1 goes first")
		conn.WriteToUDP([]byte("Player 1 goes first\n"), addr1)
		conn.WriteToUDP([]byte("Player 1 goes first\n"), addr2)
		player1.IsTurn = true
		player2.IsTurn = false
	} else {
		firstDefender = getFirstDefender(*player1Pokemons)
		fmt.Println("Player 2 goes first")
		conn.WriteToUDP([]byte("Player 2 goes first\n"), addr1)
		conn.WriteToUDP([]byte("Player 2 goes first\n"), addr2)
		player1.IsTurn = false
		player2.IsTurn = true
	}

	// The battle loop
	var player1Pokemon = firstAttacker
	var player2Pokemon = firstDefender
	for {
		if player1.IsTurn {
			if !isAlive(player1Pokemon) {
				fmt.Println(player1Pokemon.Name, "is dead")
				conn.WriteToUDP([]byte(fmt.Sprintf("%s is dead\n", player1Pokemon.Name)), addr1)
				player1Pokemon = switchPokemon(*player1Pokemons, conn, addr1)
				if player1Pokemon == nil {
					fmt.Println("Player 1 has no pokemon left")
					fmt.Println("Player 1 lost")
					conn.WriteToUDP([]byte("You have no pokemon left. You lost.\n"), addr1)
					conn.WriteToUDP([]byte("Player 1 has no pokemon left. Player 2 wins.\n"), addr2)
					break
				} else {
					fmt.Println("Player 1 switched to", player1Pokemon.Name)
					conn.WriteToUDP([]byte(fmt.Sprintf("Player 1 switched to %s\n", player1Pokemon.Name)), addr1)
				}
			}

			fmt.Printf("Player 1 turn. Your current pokemon is %s. Choose your action:\n", player1Pokemon.Name)
			conn.WriteToUDP([]byte(fmt.Sprintf("Your turn. Your current pokemon is %s. Choose your action:\n", player1Pokemon.Name)), addr1)
			command := readCommands(conn, addr1)
			switch command {
			case "attack":
				attack(player1Pokemon, player2Pokemon, conn, addr1)
			case "switch":
				displaySelectedPokemons(*player1Pokemons, conn, addr1)
				player1Pokemon = switchToChosenPokemon(*player1Pokemons, conn, addr1)
				fmt.Println("Player 1 switched to", player1Pokemon.Name)
				conn.WriteToUDP([]byte(fmt.Sprintf("Switched to %s\n", player1Pokemon.Name)), addr1)
			case "?":
				displayCommandsList(conn, addr1)
			}

			player1.IsTurn = false
			player2.IsTurn = true
		}

		if player2.IsTurn {
			if !isAlive(player2Pokemon) {
				fmt.Println(player2Pokemon.Name, "is dead")
				conn.WriteToUDP([]byte(fmt.Sprintf("%s is dead\n", player2Pokemon.Name)), addr2)
				player2Pokemon = switchPokemon(*player2Pokemons, conn, addr2)
				if player2Pokemon == nil {
					fmt.Println("Player 2 has no pokemon left")
					fmt.Println("Player 2 lost")
					conn.WriteToUDP([]byte("You have no pokemon left. You lost.\n"), addr2)
					conn.WriteToUDP([]byte("Player 2 has no pokemon left. Player 1 wins.\n"), addr1)
					break
				} else {
					fmt.Println("Player 2 switched to", player2Pokemon.Name)
					conn.WriteToUDP([]byte(fmt.Sprintf("Player 2 switched to %s\n", player2Pokemon.Name)), addr2)
				}
			}

			fmt.Printf("Player 2 turn. Your current pokemon is %s. Choose your action:\n", player2Pokemon.Name)
			conn.WriteToUDP([]byte(fmt.Sprintf("Your turn. Your current pokemon is %s. Choose your action:\n", player2Pokemon.Name)), addr2)
			command := readCommands(conn, addr2)
			switch command {
			case "attack":
				attack(player2Pokemon, player1Pokemon, conn, addr2)
			case "switch":
				displaySelectedPokemons(*player2Pokemons, conn, addr2)
				player2Pokemon = switchToChosenPokemon(*player2Pokemons, conn, addr2)
				fmt.Println("Player 2 switched to", player2Pokemon.Name)
				conn.WriteToUDP([]byte(fmt.Sprintf("Switched to %s\n", player2Pokemon.Name)), addr2)
			case "?":
				displayCommandsList(conn, addr2)
			}

			player2.IsTurn = false
			player1.IsTurn = true
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func attack(attacker *Pokemon, defender *Pokemon, conn *net.UDPConn, addr *net.UDPAddr) {
	// Calculate the damage
	var dmg float32
	var attackerMove = chooseAttack(*attacker)
	fmt.Println(attacker.Name, "chose", attackerMove.Name, "to attack", defender.Name)
	conn.WriteToUDP([]byte(fmt.Sprintf("%s chose %s to attack %s\n", attacker.Name, attackerMove.Name, defender.Name)), addr)

	switch attackerMove.Name {
	case "Tackle":
		dmg = attackerMove.Power - defender.Stats.Defense
	case "Special":
		attackingElement := attackerMove.Element
		dmgWhenAttacked := defender.DamegeWhenAttacked
		defendingElement := []string{}
		for _, element := range dmgWhenAttacked {
			defendingElement = append(defendingElement, element.Element)
		}
		highestCoefficient := float32(0)

		// Check for the highest coefficient
		for i, element := range defendingElement {
			if isContain(attackingElement, element) {
				if highestCoefficient < dmgWhenAttacked[i].Coefficient {
					highestCoefficient = dmgWhenAttacked[i].Coefficient
				}
			}
		}

		// If the attacker has an element that the defender doesn't have, set the coefficient to 1
		for _, element := range defendingElement {
			if !isContain(attackingElement, element) && highestCoefficient < 1 {
				highestCoefficient = 1
			}
		}

		dmg = attackerMove.Power*highestCoefficient - defender.Stats.Sp_Defense
	}

	if dmg < 0 {
		dmg = 0
	}
	fmt.Println(attacker.Name, "attacked", defender.Name, "with", attackerMove.Name, "and dealt", dmg, "damage")
	conn.WriteToUDP([]byte(fmt.Sprintf("%s attacked %s with %s and dealt %.2f damage\n", attacker.Name, defender.Name, attackerMove.Name, dmg)), addr)
	defender.Stats.HP -= dmg
}

func chooseAttack(pokemon Pokemon) Moves {
	n, _ := randomInt(2)
	return pokemon.Moves[n]
}

func isContain[T any](arr []T, element T) bool {
	for _, a := range arr {
		if reflect.DeepEqual(a, element) {
			return true
		}
	}
	return false
}

func getFirstAttacker(allBattlingPokemons []Pokemon) *Pokemon {
	var highestSpeed = 0
	var choosenPokemonIndex = 0
	for i, pokemon := range allBattlingPokemons {
		if pokemon.Stats.Speed > highestSpeed {
			highestSpeed = pokemon.Stats.Speed
			choosenPokemonIndex = i
		}
	}

	return &allBattlingPokemons[choosenPokemonIndex]
}

func getFirstDefender(defenderPokemons []Pokemon) *Pokemon {
	var highestSpeed = 0
	var choosenPokemonIndex = 0
	for i, pokemon := range defenderPokemons {
		if pokemon.Stats.Speed > highestSpeed {
			highestSpeed = pokemon.Stats.Speed
			choosenPokemonIndex = i
		}
	}

	return &defenderPokemons[choosenPokemonIndex]
}

func isAlive(pokemon *Pokemon) bool {
	return pokemon.Stats.HP > 0
}

func switchPokemon(pokemonsList []Pokemon, conn *net.UDPConn, addr *net.UDPAddr) *Pokemon {
	for i := 0; i < len(pokemonsList); i++ {
		if isAlive(&pokemonsList[i]) {
			return &pokemonsList[i]
		}
	}
	return nil
}

func displayCommandsList(conn *net.UDPConn, addr *net.UDPAddr) {
	conn.WriteToUDP([]byte("List of commands:\n"), addr)
	conn.WriteToUDP([]byte("\tattack: to attack the opponent\n"), addr)
	conn.WriteToUDP([]byte("\tswitch: to switch to another pokemon\n"), addr)
}

func displaySelectedPokemons(pokemonsList []Pokemon, conn *net.UDPConn, addr *net.UDPAddr) {
	conn.WriteToUDP([]byte("You have:\n"), addr)
	for i, pokemon := range pokemonsList {
		conn.WriteToUDP([]byte(fmt.Sprintf("%d. %s\n", i, pokemon.Name)), addr)
	}
	conn.WriteToUDP([]byte("Please enter the index of the pokemon you want to switch to:\n"), addr)
}

func switchToChosenPokemon(pokemonsList []Pokemon, conn *net.UDPConn, addr *net.UDPAddr) *Pokemon {
	for {
		index := readIndex(conn, addr)
		if index < 0 || index >= len(pokemonsList) {
			conn.WriteToUDP([]byte("Please enter a valid index.\n"), addr)
			continue
		}
		if isAlive(&pokemonsList[index]) {
			return &pokemonsList[index]
		} else {
			conn.WriteToUDP([]byte("This pokemon is dead. Please select another one.\n"), addr)
		}
	}
}

func readCommands(conn *net.UDPConn, addr *net.UDPAddr) string {
	buffer := make([]byte, 1024)
	n, _, _ := conn.ReadFromUDP(buffer)
	command := strings.TrimSpace(string(buffer[:n]))
	if command == "attack" || command == "switch" || command == "?" {
		return strings.ToLower(command)
	}
	conn.WriteToUDP([]byte("Please enter a valid command\n"), addr)
	return readCommands(conn, addr)
}

func readIndex(conn *net.UDPConn, addr *net.UDPAddr) int {
	buffer := make([]byte, 1024)
	n, _, _ := conn.ReadFromUDP(buffer)
	input := strings.TrimSpace(string(buffer[:n]))
	index, _ := strconv.Atoi(input)
	return index
}

func printPokemonInfo(index int, pokemon Pokemon) {
	fmt.Println(index, ":", pokemon.Name)

	fmt.Println("\tElements: ")
	for _, element := range pokemon.Elements {
		fmt.Println("\t\tElement:", element)
	}

	fmt.Println("\tStats:")
	fmt.Println("\t\tHP:", pokemon.Stats.HP)
	fmt.Println("\t\tAttack:", pokemon.Stats.Attack)
	fmt.Println("\t\tDefense:", pokemon.Stats.Defense)
	fmt.Println("\t\tSpeed:", pokemon.Stats.Speed)
	fmt.Println("\t\tSp_Attack:", pokemon.Stats.Sp_Attack)
	fmt.Println("\t\tSp_Defense:", pokemon.Stats.Sp_Defense)

	fmt.Println("\tDamage When Attacked:")
	for _, element := range pokemon.DamegeWhenAttacked {
		fmt.Printf("\t\tElement: %s. Coefficient: %f\n", element.Element, element.Coefficient)
	}
}

func selectPokemon(player *Player, conn *net.UDPConn, addr *net.UDPAddr) *[]Pokemon {
	var selectedPokemons = []Pokemon{}
	counter := 1
	for {
		if len(selectedPokemons) == 3 {
			break
		}
		conn.WriteToUDP([]byte(fmt.Sprintf("Enter the index of the %d pokemon you want to select: ", counter)), addr)
		index := readIndex(conn, addr)
		if index < 0 || index >= len(player.Inventory) {
			conn.WriteToUDP([]byte("Invalid index\n"), addr)
			continue
		}

		if isContain(selectedPokemons, player.Inventory[index]) {
			conn.WriteToUDP([]byte("You have selected this pokemon. Please select another one.\n"), addr)
			continue
		}

		conn.WriteToUDP([]byte(fmt.Sprintf("Selected %s\n", player.Inventory[index].Name)), addr)
		counter++
		selectedPokemons = append(selectedPokemons, player.Inventory[index])
	}

	conn.WriteToUDP([]byte("You have selected: "), addr)
	for _, pokemon := range selectedPokemons {
		conn.WriteToUDP([]byte(fmt.Sprintf("%s ", pokemon.Name)), addr)
	}
	conn.WriteToUDP([]byte("\n"), addr)

	return &selectedPokemons
}

// Use the assigned variables pokemon1, pokemon2, and pokemon3 as needed.

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
