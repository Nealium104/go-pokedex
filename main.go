package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Println("Welcome to the pokedex! You can ask for information on Pokemon.")
	commands := commands()
	for _, value := range commands {
		fmt.Printf("Command: %v\n", value.name)
		fmt.Printf("Description: %v\n\n", value.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Exiting...")
	os.Exit(0)
	return nil
}

// keep track of which page you're on (closure?)
// print api calls for each page you're on (from (20 - (page * 20)) to page * 20)

type Location struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

var page int

func commandMap() error {
	page++
	fmt.Printf("Locations:\n")
	fmt.Printf("Page %v\n", page)
	for i := ((page * 20) - 19); i <= (page * 20); i++ {
		res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", i))
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			log.Fatal(err)
		}
		location := Location{}
		json.Unmarshal(body, &location)
		fmt.Printf("%v\n", location.Name)
	}
	return nil
}

func commandMapb() error {
	page--
	if page < 1 {
		fmt.Println("Sorry, you've reached the beginning of the results.")
	} else {
		fmt.Printf("Locations:\n")
		fmt.Printf("Page %v\n", page)
		for i := ((page * 20) - 19); i <= (page * 20); i++ {
			res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", i))
			if err != nil {
				log.Fatal(err)
			}
			body, err := io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode > 299 {
				log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			}
			if err != nil {
				log.Fatal(err)
			}
			location := Location{}
			json.Unmarshal(body, &location)
			fmt.Printf("%v\n", location.Name)
		}
	}
	return nil
}

func cleanInput(input string) string {
	input = strings.Trim(input, " ")
	input = strings.ToLower(input)
	return input
}

func commands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent command displays the next 20 locations.",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Similar to the map command, but displays the previous 20 locations. If on the first page, it instead pritns an error",
			callback:    commandMapb,
		},
	}
}

func main() {
	fmt.Println("Welcome to the Pokedex!")
	commands := commands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nPokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		if cmd, exists := commands[command]; exists {
			cmd.callback()
		} else {
			fmt.Println("That's not a known command, please try again.")
		}
	}
}
