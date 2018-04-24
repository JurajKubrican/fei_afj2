package main

//import "./utils"

import (
	"fmt"
	"strings"
	"strconv"
)

type nState struct {
	id      string
	next    map[string][]string
	final   bool
	initial bool
}

type dState struct {
	id      map[string]struct{}
	next    map[string]string
	final   bool
	initial bool
}

func readNStates(file string) (map[string]nState, []string, string) {

	lines, _ := readLines(file)
	if lines == nil {
		fmt.Println("empty / error")
	}

	/*
	 * STATES
	 */
	initial := ""
	var states = make(map[string]nState)
	numStates, _ := strconv.Atoi(lines[0])
	var i = 2;
	for ; i < 2+numStates; i++ {
		parsedLine := strings.Split(lines[i], " ")
		newState := nState{
			id:   parsedLine[0],
			next: make(map[string][]string),
		}
		if len(parsedLine) > 1 && strings.Contains(parsedLine[1], "F") {
			newState.final = true
		}
		if len(parsedLine) > 1 && strings.Contains(parsedLine[1], "I") {
			newState.initial = true
			initial = newState.id
		}

		states[newState.id] = newState
	}

	/*
	 * Alpha
	 */
	var alphabet []string
	numChars, _ := strconv.Atoi(lines[0])
	for ; i < (2 + numStates + numChars); i++ {
		alphabet = append(alphabet, lines[i])
	}

	/*
	 * NEXT
	 */
	for ; i < len(lines); i++ {
		parsedLine := strings.Split(lines[i], ",")
		state := states[parsedLine[0]]
		next := state.next[parsedLine[1]]
		next = append(next, parsedLine[2])
		state.next[parsedLine[1]] = next
	}

	return states, alphabet, initial
}

func getDFAStateName(nStates map[string]nState, initialState string, char string) map[string]struct{} {

	nState := nStates[initialState]

	result := make(map[string]struct{})

	if val, ok := nState.next[""]; ok {
		for _, nextState := range val {
			tempName := getDFAStateName(nStates, nextState, char)
			for k, v := range tempName {
				result[k] = v
			}
		}
	}

	if val, ok := nState.next[char]; ok {
		for _, nextState := range val {
			result[nextState] = struct{}{}
		}
	}
	return result

}

//[]dState
func makeDFA(nStates map[string]nState, alpha []string, initial string) {

	initialName := getDFAStateName(nStates, initial, "a")


	stack := Stack{}
	stack.Push(initial)

	for ; stack.size > 0; {
		state:= stack.Pop()
		for _
	}

	fmt.Println(name)
	name = getDFAStateName(nStates, initial, "b")
	fmt.Println(name)

	//for _, nState := range nStates {
	//	for _, next := range nState.next {
	//		fmt.Println(next)
	//	}
	//}

}

func main() {
	states, alphabet, initial := readNStates("./in1.txt")
	makeDFA(states, alphabet, initial)
	//dStates := makeDFA(states, alphabet)

	//fmt.Println(alphabet)
	//fmt.Println(states)

}
