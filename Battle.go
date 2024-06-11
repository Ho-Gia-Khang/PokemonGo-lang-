package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	// Name              string
	// ID                string
	// PlayerCoordinateX int
	// PlayerCoordinateY int
	// Addr              *net.UDPAddr
	IsTurn			bool
	Inventory 	   []Pokemon
}

func battleScene(player1 *Player, player2 *Player) {
	
	if(len(player1.Inventory) < 3){
		fmt.Println("Player 1 has less than 3 pokemons")
		return
	} else if (len(player2.Inventory) < 3) {
		fmt.Println("Player 2 has less than 3 pokemons")
		return
	}

	// init the battle
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Battle Start!")

	// player 1 select 3 pokemons
	fmt.Println("Player 1 please select 3 pokemons from: ")
	for i := range(len(player1.Inventory)){
		printPokemonInfo(i, player1.Inventory[i])
	}
	player1Pokemons := selectPokemon(player1, reader)

	// player 2 select 3 pokemons
	fmt.Println("Player 2 please select 3 pokemons from: ")
	for i := range(len(player2.Inventory)){
		printPokemonInfo(i, player2.Inventory[i])
	}
	player2Pokemons := selectPokemon(player2, reader)

	var allBattlingPokemons = append(*player1Pokemons, *player2Pokemons...)
	var firstAttacker = getFirstAttacker(allBattlingPokemons)
	var firstDefender *Pokemon

	if(isContain(*player1Pokemons, *firstAttacker)){
		firstDefender = getFirstDefender(*player2Pokemons)
		player1.IsTurn = false
		player2.IsTurn = true
	} else {
		firstDefender = getFirstDefender(*player1Pokemons)
		player1.IsTurn = true
		player2.IsTurn = false
	}

	attack(firstAttacker, firstDefender)
	var player1Pokemon = firstAttacker
	var player2Pokemon = firstDefender

	// the battle loop
	for {
		if isLost(player1Pokemons) {
			fmt.Println("Player 1 lost")
			break
		} else if isLost(player2Pokemons) {
			fmt.Println("Player 2 lost")
			break
		}

		if player1.IsTurn {
			fmt.Print("Player 1 turn. Your current pokemon is ", player1Pokemon.Name, ". Choose your action:\n")
			command := readCommands(reader)
			
			if isAlive(player1Pokemon){
				continue
			} else {
				fmt.Println(player1Pokemon.Name, "is dead")
				player1Pokemon = switchPokemon(*player1Pokemons)
				if player1Pokemon == nil {
					fmt.Println("Player 1 has no pokemon left")
				} else {
					fmt.Println("Player 1 switched to", player1Pokemon.Name)
				}
			}
			switch command {
				case "attack":
					attack(player1Pokemon, player2Pokemon)
				case "switch":
					player1Pokemon = switchToChosenPokemon(*player1Pokemons, readIndex(reader))
			}
			
			player1.IsTurn = false
			player2.IsTurn = true
		}
		
		if player2.IsTurn {
			fmt.Print("Player 1 turn. Your current pokemon is ", player2Pokemon.Name, ". Choose your action:\n")
			command := readCommands(reader)
			
			if isAlive(player2Pokemon){
				continue
			} else {
				fmt.Println(player2Pokemon.Name, "is dead")
				player2Pokemon = switchPokemon(*player2Pokemons)
				if player2Pokemon == nil {
					fmt.Println("Player 1 has no pokemon left")
				} else {
					fmt.Println("Player 1 switched to", player2Pokemon.Name)
				}
			}
			switch command {
				case "attack":
					attack(player2Pokemon, player2Pokemon)
				case "switch":
					player2Pokemon = switchToChosenPokemon(*player2Pokemons, readIndex(reader))
			}

			player2.IsTurn = false
			player1.IsTurn = true
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func attack(attacker *Pokemon, defender *Pokemon) {
	// Calculate the damage
	var dmg float32
	var attackerMove = chooseAttack(*attacker)
	fmt.Println(attacker.Name, "chosed", attackerMove.Name, "to attack", defender.Name)

	switch attackerMove.Name{
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

			// check for the highest coefficient
			for i, element := range defendingElement {
				if isContain(attackingElement, element) {
					if highestCoefficient < dmgWhenAttacked[i].Coefficient {
						highestCoefficient = dmgWhenAttacked[i].Coefficient
					}
				}
			}

			// if the attacker have the element that the defender doesn't have, set the coefficient to 1
			for _, element := range defendingElement {
				if !isContain(attackingElement, element) && highestCoefficient < 1{
					highestCoefficient = 1
				}
			}

			dmg = attackerMove.Power * highestCoefficient - defender.Stats.Sp_Defense
	}
	
	if (dmg < 0){
		dmg = 0
	}
	fmt.Println(attacker.Name, "attacked", defender.Name, "with", attackerMove.Name, "and dealt", dmg, "damage")
	defender.Stats.HP -= dmg
}

func chooseAttack(pokemon Pokemon) Moves {
	n := rand.Intn(2)
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
	for i, pokemon := range allBattlingPokemons{
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
	for i, pokemon := range defenderPokemons{
		if pokemon.Stats.Speed > highestSpeed {
			highestSpeed = pokemon.Stats.Speed
			choosenPokemonIndex = i
		}
	}

	return &defenderPokemons[choosenPokemonIndex]
}

func isLost(pokemonsList *[]Pokemon) bool {
	for _, pokemon := range *pokemonsList {
		if pokemon.Stats.HP > 0 {
			return false
		}
	}
	return true
}

func isAlive(pokemon *Pokemon) bool {
	return pokemon.Stats.HP > 0
}

func switchPokemon(pokemonsList []Pokemon) *Pokemon {
	if isAlive(&pokemonsList[0]){
		return &pokemonsList[0]
	} else if isAlive(&pokemonsList[1]){
		return &pokemonsList[1]
	} else if isAlive(&pokemonsList[2]) {
		return &pokemonsList[2]
	} else {
		return nil
	}
}

func switchToChosenPokemon(pokemonsList []Pokemon, index int) *Pokemon {
	fmt.Print("You have: ")
	for _, pokemon := range pokemonsList{
		if isAlive(&pokemon){
			fmt.Print(pokemon.Name, " ")
		}
	}
	for{
		if isAlive(&pokemonsList[index]){
			return &pokemonsList[index]
		} else {
			fmt.Println("This pokemon is dead")
		}
	}	
}

func readCommands(reader *bufio.Reader) string {
	// read the commands from the user
	for{
		input, _ := reader.ReadString('\n')
		command := strings.TrimSpace(input)
		commands := strings.Split(command, " ")
		if (len(commands) > 1){
			fmt.Println("Please enter a command with one word")
		} else if (commands[0] == "attack" || commands[0] == "switch"){
			return strings.ToLower(commands[0]) 
		} else {
			fmt.Println("Please enter a valid command")
		}
	}
}

func readIndex(reader *bufio.Reader) int {
	// read the index from the user
	input, _ := reader.ReadString('\n')
	index, _ := strconv.Atoi(strings.TrimSpace(input))
	return index
}	

func printPokemonInfo(index int, pokemon Pokemon){
	fmt.Println(index, ":", pokemon.Name)

	fmt.Println("\tElements: ")
	for _, element := range pokemon.Elements{
		fmt.Println("\t\tElement:", element)
	}

	fmt.Println("\tStats:")
	fmt.Println("\t\tHP:", pokemon.Stats.HP)
	fmt.Println("\t\tAttack:", pokemon.Stats.Attack)
	fmt.Println("\t\tDefense:", pokemon.Stats.Defense)
	fmt.Println("\t\tSpeed:", pokemon.Stats.Speed)
	fmt.Println("\t\tSp_Attack:", pokemon.Stats.Sp_Attack)
	fmt.Println("\t\tSp_Defense:", pokemon.Stats.Sp_Defense)

	fmt.Println("\tDamege When Attacked:")
	for _, element := range pokemon.DamegeWhenAttacked{
		fmt.Printf("\t\tElement: %s. Coefficient: %f\n", element.Element ,element.Coefficient)
	}
}

func selectPokemon(player *Player, reader *bufio.Reader) *[]Pokemon {
	var selectedPokemons = []Pokemon{}
	counter := 1
	for {
		if len(selectedPokemons) == 3 {
			break
		}
		fmt.Printf("Enter the index of the %d you want to select: ", counter)
		index := readIndex(reader)
		if index < 0 || index >= len(player.Inventory) {
			fmt.Println("Invalid index")
			continue
		}

		if isContain(selectedPokemons, player.Inventory[index]) {
			fmt.Println("You have selected this pokemon. Please select another one.")
			continue
		}

		fmt.Println("Selected", player.Inventory[index].Name)
		counter++
		selectedPokemons = append(selectedPokemons, player.Inventory[index])
	}

	fmt.Println("You have selected: ")
	for _, pokemon := range selectedPokemons {
		fmt.Print(pokemon.Name, " ")
	}
	fmt.Println()

	return &selectedPokemons
}