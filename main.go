package main

import (
	"bufio"
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

func commandMap() error {
	res, err := http.Get("https://pokeapi.co/api/v2/location/1")
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
	fmt.Printf("%s", body)
	return nil
}

func commandMapb() error {
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
	print("Welcome to the Pokedex!\n")
	commands := commands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		command := cleanInput(scanner.Text())
		if cmd, exists := commands[command]; exists {
			cmd.callback()
		} else {
			fmt.Println("That's not a known command, please try again.")
		}
	}
}
