package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Stats struct {
	HP                 int `json:"HP"`
	Attack 		   		int `json:"Attack"`                  
	Defense            int `json:"Defense"`                 
	Speed              int `json:"Speed"`                  
	Sp_Attack          int `json:"Sp_Attack"`                  
	Sp_Defense         int `json:"Sp_Defense"`                  
}

type GenderRatio struct {
	MaleRatio  float32 `json:"MaleRatio"` 
	FemaleRatio float32 `json:"FemaleRatio"` 
}

type Profile struct {
	Name 			string `json:"Name"`
	Height          float32 `json:"Height"`               
	Weight          float32 `json:"Weight"`              
	CatchRate 	 float32 `json:"CatchRate"`                  
	GenderRatio     GenderRatio `json:"GenderRatio"` 
	EggGroup        string `json:"EggGroup"` 
	HatchSteps	  int `json:"HatchSteps"`
	Abilities	   string `json:"Abilities"` 
}

type DamegeWhenAttacked struct {
	Element     string `json:"Element"`  
	Coefficient float32 `json:"Coefficient"` 
}

type Moves struct {
	Name string `json:"Name"`
	Element string `json:"Element"`
	Power string `json:"Power"`
	Acc 	int `json:"Acc"`
	PP 		int `json:"PP"`
	Description string `json:"Description"`
}

type Pokemon struct {
	Elements           []string `json:"Elements"`            
	EV                 int `json:"EV"`
	Stats              Stats `json:"Stats"`                                
	Profile			Profile `json:"Profile"`              
	DamegeWhenAttacked []DamegeWhenAttacked `json:"DamegeWhenAttacked"`
	EvolutionLevel    int `json:"EvolutionLevel"`
	NextEvolution     string `json:"NextEvolution"`
	Moves			  []Moves `json:"Moves"`                  
}

const (
	numberOfPokemons = 2
)

var pokemons []Pokemon;

func main(){
	c := colly.NewCollector()
	
	for i := 1; i <= numberOfPokemons; i++{
		url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", i)
		crawlPokemons(c, url)
		refresh(url)
	}
}

func crawlPokemons(c *colly.Collector, url string,){
	c.OnRequest(func(r *colly.Request) {
    	fmt.Println("Visiting:", r.URL.String())
	})

	c.OnHTML("div.mui-panel", func(e *colly.HTMLElement) {
		stats := Stats{}
		e.ForEach("div.detail-panel-content > div.detail-header > div.detail-infobox > div.detail-stats > div.detail-stats-row", func(_ int, el *colly.HTMLElement) {
			title := el.ChildText("span:not([class])")
			switch title {
				case "HP":
					stats.HP, _ = strconv.Atoi(el.ChildText("span.stat-bar > div.stat-bar-fg"));
				case "Attack":
					stats.Attack, _ = strconv.Atoi(el.ChildText("span.stat-bar > div.stat-bar-fg"));
				case "Defense":
					stats.Defense, _ = strconv.Atoi(el.ChildText("span.stat-bar > div.stat-bar-fg"));
				case "Speed":
					stats.Speed, _ = strconv.Atoi(el.ChildText("span.stat-bar > div.stat-bar-fg"));
				case "Sp Atk":
					stats.Sp_Attack, _ = strconv.Atoi(el.ChildText("span.stat-bar > div.stat-bar-fg"));
				case "Sp Def":
					stats.Sp_Defense, _ = strconv.Atoi(el.ChildText("span.stat-bar > div.stat-bar-fg"));
				default:
					fmt.Println("Unknown title: ", title)
			}
		})

		genderRatio := GenderRatio{}
		profile := Profile{}
		profile.Name = e.ChildText("h1.detail-panel-header")
		e.ForEach("div.detail-below-header > div.monster-minutia", func(_ int, el *colly.HTMLElement) {
			title1 := el.ChildText("strong:nth-child(1)")
			detail1 := el.ChildText("span:nth-child(2)")
			switch title1 {
				case "Height:":
					height := strings.Split(detail1, " ")
					val, _ := strconv.ParseFloat(height[0], 32)
					profile.Height = float32(val)
				case "Catch Rate:":
					rates := strings.Split(detail1, "")
					val, _ := strconv.ParseFloat(rates[0], 32)
					profile.CatchRate = float32(val)
				case "Egg Groups:":
					profile.EggGroup = detail1
				case "Abilities:":
					profile.Abilities = detail1
				default:
					fmt.Println("Unknown title: ", title1)
			}

			title2 := el.ChildText("strong:nth-child(3)")
			detail2 := el.ChildText("span:nth-child(4)")

			switch title2 {
				case "Weight:":
					weight := strings.Split(detail1, " ")
					val, _ := strconv.ParseFloat(weight[0], 32)
					profile.Weight = float32(val)
				case "Gender Ratio:":
					ratios := strings.Split(detail2, " ")

					ratio1 := strings.Split(ratios[0], "%")
					val1, _ := strconv.ParseFloat(ratio1[0], 32)
					genderRatio.MaleRatio = float32(val1)

					ratio2 := strings.Split(ratios[2], "%")
					val2, _ := strconv.ParseFloat(ratio2[0], 32)
					genderRatio.FemaleRatio = float32(val2)
					profile.GenderRatio = genderRatio
				case "Hatch Steps:":
					profile.HatchSteps, _ = strconv.Atoi(detail2)
			}
		})
		profile.GenderRatio = genderRatio

		damegeWhenAttacked := []DamegeWhenAttacked{}
		e.ForEach("div.detail-below-header > div.when-attacked > div.when-attacked-row", func(_ int, el *colly.HTMLElement) {
			monsterType1 := el.ChildText("span.monster-type:nth-child(1)")
			multiplier1 := el.ChildText("span.monster-multiplier:nth-child(2)")
			multipliers1 := strings.Split(multiplier1, "x")
			coefficient1, _ := strconv.ParseFloat(multipliers1[0], 32)
			
			damegeWhenAttacked = append(damegeWhenAttacked, DamegeWhenAttacked{
				Element: monsterType1,
				Coefficient: float32(coefficient1),
			})

			monsterType2 := el.ChildText("span.monster-type:nth-child(3)")
			multiplier2 := el.ChildText("span.monster-multiplier:nth-child(4)")
			multipliers2 := strings.Split(multiplier2, "x")
			coefficient2, _ := strconv.ParseFloat(multipliers2[0], 32)
			
			damegeWhenAttacked = append(damegeWhenAttacked, DamegeWhenAttacked{
				Element: monsterType2,
				Coefficient: float32(coefficient2),
			})
		})

		moves := []Moves{}
		e.ForEach("div.detail-below-header > div.monster-moves > div.moves-row", func(_ int, el *colly.HTMLElement) {
			fmt.Println("Moves: ")
		});
		
		pokemon := Pokemon{
			Stats: stats,
			Profile: profile,
			DamegeWhenAttacked: damegeWhenAttacked,
			Moves: moves,
		}

		e.ForEach("div.detail-header > div.detail-infobox > div.detail-types-and-num > div.detail-types > span.monster-type", func(_ int, el *colly.HTMLElement) {
			pokemon.Elements = append(pokemon.Elements, el.Text)
		})

		evolutionLabel := e.ChildText("div.detail-below-header > div.evolutions > div.evolution-row > div.evolution-label > span")
		labels := strings.Split(evolutionLabel, " ")
		nextEvolution := labels[3]
		evolutionStage := labels[len(labels)-1]
		evolutionStage = strings.Replace(evolutionStage, ".", "", -1)
		evolutionLevel, _ := strconv.Atoi(evolutionStage)
		pokemon.EvolutionLevel = evolutionLevel
		pokemon.NextEvolution = nextEvolution

		js, err := json.MarshalIndent(pokemon, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))

		pokemons = append(pokemons, pokemon)
		
	})

	c.Visit(url)
}

func refresh(url string){
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the page:", err)
	}

	// _ , err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading the response body:", err)
	// 	resp.Body.Close()
	// }

	//fmt.Println("Page content:", string(body))
	resp.Body.Close()
}