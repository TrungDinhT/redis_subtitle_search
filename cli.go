package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/mattn/go-tty"
)

var BACKSPACE_ONE_CHAR = []byte("\b \b")

var MOVE_CURSOR_NEXT_LINE = []byte("\x1b[E")     //begin of next line
var MOVE_CURSOR_PREVIOUS_LINE = []byte("\x1b[F") //begin of previous line

var DELETE_WHOLE_LINE = []byte("\x1b[0K")

var MOVE_CURSOR_UP = []byte("\x1b[A")
var MOVE_CURSOR_DOWN = []byte("\x1b[B")
var MOVE_CURSOR_RIGHT = []byte("\x1b[C")
var MOVE_CURSOR_LEFT = []byte("\x1b[D")

var CARRIAGE_RETURN = "\r"
var NEW_LINE = "\n"

func display(tty *tty.TTY, resultsChan chan []string, phraseChan chan []rune) {
	var results []string
loop:
	for {
		phrase := <-phraseChan

		if len(phrase) == 1 && (phrase[0] == 13 || phrase[0] == 27) {
			break loop
		}

		prevResults := results
		results = <-resultsChan

		i := 0
		for ; i < len(results); i++ {
			tty.Output().Write([]byte(NEW_LINE + CARRIAGE_RETURN + results[i]))
			tty.Output().Write(DELETE_WHOLE_LINE)
		}
		for ; i < len(prevResults); i++ {
			tty.Output().Write(MOVE_CURSOR_NEXT_LINE)
			tty.Output().Write(DELETE_WHOLE_LINE)
		}
		for ; i > 0; i-- {
			tty.Output().Write(MOVE_CURSOR_UP)
		}

		tty.Output().Write([]byte(CARRIAGE_RETURN + string(phrase)))
		tty.Output().Write(DELETE_WHOLE_LINE)
	}
}

func run(tty *tty.TTY, resultsChan chan []string, phraseChan chan []rune) {
	var phrase []rune
loop:
	for {
		r, err := tty.ReadRune()

		if err != nil {
			log.Fatal(err)
		}

		switch r {
		case 13, 27:
			phraseChan <- []rune{r}
			resultsChan <- make([]string, 0)
			break loop
		case 8, 127:
			if len(phrase) > 0 {
				phrase = phrase[:len(phrase)-1]
			}
		default:
			phrase = append(phrase, r)
		}
		phraseChan <- phrase
		resultsChan <- randomResult(phrase)

	}
}

// func searchResult(phraseChan chan []rune, resultsChan chan []string) {
// 	for {
// 		phrase := <-phraseChan
// 		if len(phrase) == 0 {
// 			break
// 		}
// 		resultsChan <- randomResult(phrase)
// 	}
// }

func randomResult(phrase []rune) []string {
	if len(phrase) == 0 {
		return make([]string, 0)
	}

	n := rand.Intn(5) + 1
	results := make([]string, 0, n)
	for i := 0; i < n; i++ {
		result := ""
		for j := 0; j < rand.Intn(10)+1; j++ {
			result += strconv.Itoa(i + 1)
		}
		results = append(results, result)
	}

	return results
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	tty, err := tty.Open()

	if err != nil {
		log.Fatal(err)
	}

	phraseChan := make(chan []rune, 10)
	resultsChan := make(chan []string, 10)

	fmt.Println("Run")

	go run(tty, resultsChan, phraseChan)

	display(tty, resultsChan, phraseChan)
}
