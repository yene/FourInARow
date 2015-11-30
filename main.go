package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Documentation can be found on http://theaigames.com/competitions/four-in-a-row/getting-started

type Settings struct {
	timebank      int      // Maximum time in milliseconds that your bot can have in its time bank.
	time_per_move int      // Time in milliseconds that is added to your bot's time bank each move.
	player_names  []string // A list of all player names in this match, including your bot's name.
	your_bot      string   // The name of your bot for this match.
	your_botid    int      // The number used in a field update as your bot's chips.
	field_columns int      // The number of columns of the playing field.
	field_rows    int      // The number of rows of the playing field.
}

type Match struct {
	round int     // The number of the current round.
	field string  // The complete playing field in the current game state
	grid  [][]int // Converted field to a grid.
}

var settings Settings
var match Match

var debug bool = true

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	reader := bufio.NewReader(os.Stdin)

	if debug { // use for debugging
		f, _ := os.Open("example-setup.txt")
		reader = bufio.NewReader(f)
	}

	writer := bufio.NewWriter(os.Stdout)
	logger := bufio.NewWriter(os.Stderr)

	for {
		engineOutput, isPrefix, _ := reader.ReadLine()
		if isPrefix {
			panic("Line too long for buffer")
		}

		if len(engineOutput) == 0 {
			continue
		}
		words := strings.Split(string(engineOutput), " ")

		switch words[0] {

		case "settings":
			switch words[1] {
			case "timebank":
				settings.timebank, _ = strconv.Atoi(words[2])
			case "time_per_move":
				settings.time_per_move, _ = strconv.Atoi(words[2])
			case "player_names":
				settings.player_names = words[2:3]
			case "your_bot":
				settings.your_bot = words[2]
			case "your_botid":
				settings.your_botid, _ = strconv.Atoi(words[2])
			case "field_columns":
				settings.field_columns, _ = strconv.Atoi(words[2])
			case "field_rows":
				settings.field_rows, _ = strconv.Atoi(words[2])
			default:
				logger.WriteString("Settings parameter " + words[1] + " not documented \n")
				logger.Flush()
			}

		case "update":
			switch words[2] {
			case "round":
				match.round, _ = strconv.Atoi(words[3])
			case "field":
				match.field = words[3]
				generateGrid()
				if debug {
					printGrid()
				}
			}

		case "action":
			writer.WriteString(turn() + "\n")
			writer.Flush()

		default:
			logger.WriteString("Command " + words[0] + " not documented \n")
			logger.Flush()
		}
	}
}

func generateGrid() {
	match.grid = make([][]int, settings.field_rows)
	lines := strings.Split(match.field, ";")
	for x, line := range lines {
		match.grid[x] = make([]int, settings.field_columns)
		rows := strings.Split(line, ",")
		for y, row := range rows {
			match.grid[x][y], _ = strconv.Atoi(row)
		}
	}
}

func printGrid() {
	fmt.Println("--------------")
	fmt.Println("Round ", match.round)
	fmt.Println("")
	for _, line := range match.grid {
		for _, row := range line {
			switch row {
			case 0:
				fmt.Print("⚪️ ")
			case 1:
				fmt.Print("🔴 ")
			case 2:
				fmt.Print("🔵 ")
			}
		}
		fmt.Println("")
	}
}

func turn() string {
	// it is actually grid[y][x]
	// it goes from the top left down

	// place mid until mid is full
	if match.grid[0][3] == 0 {
		return "place_disc 3"
	}

	// if we have someone 3 diagonal, horizontal or vertical finish the game

	// if enemy has somewhere 3 prevent it

	// if we have somewhere more than 2 vertical attach to them

	// else random
	return fmt.Sprintf("place_disc %d", rand.Intn(settings.field_columns))
}
