package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strings"
)

type Player struct {
	Name              string
	ID                string
	PlayerCoordinateX int
	PlayerCoordinateY int
	Addr              *net.UDPAddr
}

var players = make(map[string]*Player)
var Pokeworld [1000][1000]string

// func generateUniqueIDInt() int {
// 	var id int
// 	for i := 1; i < 1000000; i++ {
// 		if players[i] == nil {
// 			id = i
// 			return id
// 		}
// 	}
// 	return id
// }

func PokeCat(Id string, playername string, x int, y int) {
	// Check if the coordinates are within the bounds of Pokeworld.
	if x >= 0 && x < 1000 && y >= 0 && y < 1000 {
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
			fmt.Println("Player placed at", x, y)
		} else {
			// Handle the case where the position is already occupied.
			fmt.Println("Position is already occupied.")
		}
	} else {
		// Handle the case where the coordinates are out of bounds.
		fmt.Println("Coordinates are out of bounds.")
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
			xBigInt, _ := rand.Int(rand.Reader, big.NewInt(1000))
			yBigInt, _ := rand.Int(rand.Reader, big.NewInt(1000))
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

		// if parts[0] == "PokeCat" {
		// 	xBigInt, _ := rand.Int(rand.Reader, big.NewInt(1000))
		// 	yBigInt, _ := rand.Int(rand.Reader, big.NewInt(1000))
		// 	x := int(xBigInt.Int64())
		// 	y := int(yBigInt.Int64())
		// 	PokeCat(idStr, "Tin", x, y)
		// 	fmt.Println(idStr, "placed at", x, y)
		// }
		// if parts[0] == "PokeBat" {
		// 	PokeBat()

	} // }
}
