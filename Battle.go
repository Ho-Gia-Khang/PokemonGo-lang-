package main

import "math/rand"

func chooseAttack(pokemon Pokemon) Moves {
	n := rand.Intn(2)
	return pokemon.Moves[n]
}

func attack(attacker *Pokemon, defender *Pokemon) {
	// Calculate the damage
	var dmg float32
	var attackerMove = chooseAttack(*attacker)

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
			highestCoefficient := float32(1)

			for i, element := range defendingElement {
				if isContain(attackingElement, element) {
					if highestCoefficient < dmgWhenAttacked[i].Coefficient {
						highestCoefficient = dmgWhenAttacked[i].Coefficient
					}
				}
			}

			dmg = attackerMove.Power * highestCoefficient - defender.Stats.Sp_Defense
	}

	defender.Stats.HP -= dmg
}

func isContain(arr []string, element string) bool {
	for _, a := range arr {
		if a == element {
			return true
		}
	}
	return false
}