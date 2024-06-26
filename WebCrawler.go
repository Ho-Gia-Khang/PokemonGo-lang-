package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

const (
	numberOfPokemons = 649
	baseURL          = "https://pokedex.org/#/"
)

var pokemons = []Pokemon{}

func crawlPokemonsDriver(numsOfPokemons int) {

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	page.Goto(baseURL)

	for i := range numsOfPokemons {
		// simulate clicking the button to open the pokemon details
		page.Reload()
		page.WaitForURL(baseURL)

		locator := fmt.Sprintf("button.sprite-%d", i+1)
		button := page.Locator(locator).First()
		time.Sleep(500 * time.Millisecond)
		button.Click()

		url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i+1)
		page.WaitForURL(url)

		fmt.Print("Pokemon ", i+1, " ")
		newPokemon := crawlPokemons(page)
		createMoves(&newPokemon)
		pokemons = append(pokemons, newPokemon)

		page.Goto(baseURL)
	}

	// parse the pokemons variable to json file
	dataPokemons, _ := json.MarshalIndent(pokemons, "", " ")
	os.WriteFile("pokedex.json", dataPokemons, 0666)

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

func crawlPokemons(page playwright.Page) Pokemon {
	pokemon := Pokemon{}

	stats := Stats{}
	entries, _ := page.Locator("div.detail-panel-content > div.detail-header > div.detail-infobox > div.detail-stats > div.detail-stats-row").All()
	for _, entry := range entries {
		title, _ := entry.Locator("span:not([class])").TextContent()
		switch title {
		case "HP":
			stat, _ := entry.Locator("span.stat-bar > div.stat-bar-fg").TextContent()
			hp, _ := strconv.ParseFloat(stat, 32)
			stats.HP = float32(hp)
		case "Attack":
			stat, _ := entry.Locator("span.stat-bar > div.stat-bar-fg").TextContent()
			Attack, _ := strconv.ParseFloat(stat, 32)
			stats.Attack = float32(Attack)
		case "Defense":
			stat, _ := entry.Locator("span.stat-bar > div.stat-bar-fg").TextContent()
			Defense, _ := strconv.ParseFloat(stat, 32)
			stats.Defense = float32(Defense)
		case "Speed":
			speed, _ := entry.Locator("span.stat-bar > div.stat-bar-fg").TextContent()
			stats.Speed, _ = strconv.Atoi(speed)
		case "Sp Atk":
			sp_Attack, _ := entry.Locator("span.stat-bar > div.stat-bar-fg").TextContent()
			sp_attacks, _ := strconv.ParseFloat(sp_Attack, 32)
			stats.Sp_Attack = float32(sp_attacks)
		case "Sp Def":
			sp_Defense, _ := entry.Locator("span.stat-bar > div.stat-bar-fg").TextContent()
			sp_defense, _ := strconv.ParseFloat(sp_Defense, 32)
			stats.Sp_Defense = float32(sp_defense)
		default:
			fmt.Println("Unknown title: ", title)
		}
	}
	pokemon.Stats = stats

	name, _ := page.Locator("div.detail-panel > h1.detail-panel-header").TextContent()
	pokemon.Name = name

	genderRatio := GenderRatio{}
	profile := Profile{}
	entries, _ = page.Locator("div.detail-panel-content > div.detail-below-header > div.monster-minutia").All()
	for _, entry := range entries {
		title1, _ := entry.Locator("strong:not([class]):nth-child(1)").TextContent()
		stat1, _ := entry.Locator("span:not([class]):nth-child(2)").TextContent()
		switch title1 {
		case "Height:":
			heights := strings.Split(stat1, " ")
			height, _ := strconv.ParseFloat(heights[0], 32)
			profile.Height = float32(height)
		case "Catch Rate:":
			catchRates := strings.Split(stat1, "%")
			catchRate, _ := strconv.ParseFloat(catchRates[0], 32)
			profile.CatchRate = float32(catchRate)
		case "Egg Groups:":
			profile.EggGroup = stat1
		case "Abilities:":
			profile.Abilities = stat1
		}

		title2, _ := entry.Locator("strong:not([class]):nth-child(3)").TextContent()
		stat2, _ := entry.Locator("span:not([class]):nth-child(4)").TextContent()
		switch title2 {
		case "Weight:":
			weights := strings.Split(stat2, " ")
			weight, _ := strconv.ParseFloat(weights[0], 32)
			profile.Weight = float32(weight)
		case "Gender Ratio:":
			if stat2 == "N/A" {
				genderRatio.MaleRatio = 0
				genderRatio.FemaleRatio = 0
			} else {
				ratios := strings.Split(stat2, " ")

				maleRatios := strings.Split(ratios[0], "%")
				maleRatio, _ := strconv.ParseFloat(maleRatios[0], 32)
				genderRatio.MaleRatio = float32(maleRatio)

				femaleRatios := strings.Split(ratios[2], "%")
				femaleRatio, _ := strconv.ParseFloat(femaleRatios[0], 32)
				genderRatio.FemaleRatio = float32(femaleRatio)
			}

			profile.GenderRatio = genderRatio
		case "Hatch Steps:":
			profile.HatchSteps, _ = strconv.Atoi(stat2)
		}
	}
	pokemon.Profile = profile

	damegeWhenAttacked := []DamegeWhenAttacked{}
	entries, _ = page.Locator("div.when-attacked > div.when-attacked-row").All()
	for _, entry := range entries {
		element1, _ := entry.Locator("span.monster-type:nth-child(1)").TextContent()
		coefficient1, _ := entry.Locator("span.monster-multiplier:nth-child(2)").TextContent()
		coefficients1 := strings.Split(coefficient1, "x")
		coef1, _ := strconv.ParseFloat(coefficients1[0], 32)

		element2, _ := entry.Locator("span.monster-type:nth-child(3)").TextContent()
		coefficient2, _ := entry.Locator("span.monster-multiplier:nth-child(4)").TextContent()
		coefficients2 := strings.Split(coefficient2, "x")
		coef2, _ := strconv.ParseFloat(coefficients2[0], 32)

		damegeWhenAttacked = append(damegeWhenAttacked, DamegeWhenAttacked{Element: element1, Coefficient: float32(coef1)})
		damegeWhenAttacked = append(damegeWhenAttacked, DamegeWhenAttacked{Element: element2, Coefficient: float32(coef2)})
	}
	pokemon.DamegeWhenAttacked = damegeWhenAttacked

	entries, _ = page.Locator("div.evolutions > div.evolution-row").All()
	for _, entry := range entries {
		evolutionLabel, _ := entry.Locator("div.evolution-label > span").TextContent()
		evolutionLabels := strings.Split(evolutionLabel, " ")

		if evolutionLabels[0] == name {
			evolutionLevels := strings.Split(evolutionLabels[len(evolutionLabels)-1], ".")
			evolutionLevel, _ := strconv.Atoi(evolutionLevels[0])
			pokemon.EvolutionLevel = evolutionLevel

			nextEvolution := evolutionLabels[3]
			pokemon.NextEvolution = nextEvolution
		}
	}

	entries, _ = page.Locator("div.detail-types > span.monster-type").All()
	for _, entry := range entries {
		element, _ := entry.TextContent()
		pokemon.Elements = append(pokemon.Elements, element)
	}

	fmt.Println(name)
	return pokemon
}

func createMoves(pokemon *Pokemon) {
	normalMove := Moves{Name: "Tackle", Element: []string{""}, Power: pokemon.Stats.Attack, Acc: 100, PP: 35}
	specialMove := Moves{Name: "Special", Element: pokemon.Elements, Power: pokemon.Stats.Sp_Attack, Acc: 100, PP: 25}
	pokemon.Moves = append(pokemon.Moves, normalMove)
	pokemon.Moves = append(pokemon.Moves, specialMove)
}