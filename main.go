package main

import (
	"os"
	"bufio"
	"fmt"
	"strings"
	"strconv"
)

type state struct {
	id      string
	arc     []string
	final   bool
	initial bool
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func main() {
	lines, _ := readLines("./in1.txt")
	if lines != nil {
		fmt.Println("empty / error")
	}

	var initialState string
	var states = make(map[string]state)

	numStates, _ := strconv.Atoi(lines[0])
	for i := 2; i < 2+numStates; i++ {
		parsedLine := strings.Split(lines[i], " ")
		newState := state{
			id: parsedLine[0],
		}
		if len(parsedLine) > 1 && strings.Contains(parsedLine[1], "F") {
			newState.final = true
		}
		if len(parsedLine) > 1 && strings.Contains(parsedLine[1], "I") {
			initialState = newState.id
		}

		states[newState.id] = newState
	}

	fmt.Println(initialState)
	fmt.Println(states)

	var alphabet []string
	numChars, _ := strconv.Atoi(lines[0])
	for i := 2 + numStates; i < (2 + numStates + numChars); i++ {
		alphabet = append(alphabet, lines[i])
	}
	fmt.Println(alphabet)

}
