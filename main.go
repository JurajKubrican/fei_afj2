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
	next    map[string]map[string]struct{}
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
	i := 2
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
	numChars, _ := strconv.Atoi(lines[1])
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

func writeDStates(dStates []dState, alphabet []string, fileName string) {

	result := make([]string, 2)
	result[0] = strconv.Itoa(len(dStates))
	result[1] = strconv.Itoa(len(alphabet))

	nexts := make([]string, 0)
	for _, state := range dStates {
		stateLine := strMapToStr(state.id)
		if state.initial || state.final {
			stateLine = stateLine + " "
			if state.initial {
				stateLine = stateLine + "I"
			}
			if state.final {
				stateLine = stateLine + "F"
			}
		}

		result = append(result, stateLine)
		for char, next := range state.next {
			nextStr := strMapToStr(state.id) + "," + char + "," + strMapToStr(next)
			nexts = append(nexts, nextStr)
		}

	}

	for _, char := range alphabet {
		result = append(result, char)
	}

	result = append(result, nexts...)

	writeLines(result, fileName)
}

func getDFANextIDs(origState nState, result map[string]map[string]struct{}) map[string]map[string]struct{} {
	for char, states := range origState.next {
		if char == "" {
			continue
		}
		if _, ok := result[char]; !ok {
			result[char] = make(map[string]struct{})
		}
		for _, state := range states {
			result[char][state] = struct{}{}
		}
	}
	return result
}

func getDfaStateClosure(nStates map[string]nState, startIDS map[string]struct{}) dState {

	resultState := dState{
		id:   startIDS,
		next: make(map[string]map[string]struct{}),
	}
	// FOR  ID IN multiple Nstates
	for startID := range startIDS {
		nState := nStates[startID]
		resultState.final = resultState.final || nState.final      //if one of them is final
		resultState.next = getDFANextIDs(nState, resultState.next) // append all next states

		//ONLY EPS CLOSURES
		if val, ok := nState.next[""]; ok {
			for _, closureState := range val {

				nextStateID := make(map[string]struct{})
				nextStateID[closureState] = struct{}{}

				//GET NEXT IDS
				tempState := getDfaStateClosure(nStates, nextStateID)
				for key := range tempState.id {
					resultState.id[key] = struct{}{}
					resultState.final = resultState.final || tempState.final
					resultState.next = getDFANextIDs(nStates[closureState], resultState.next)
				}
			}
		}

	}

	return resultState

}

//[]dState
func makeDFA(nStates map[string]nState, alpha []string, initial string) ([]dState, []string) {

	dStates := make(map[string]dState)

	// epsilony na vstupnom stave
	initialMap := make(map[string]struct{})
	initialMap[initial] = struct{}{}
	initialState := getDfaStateClosure(nStates, initialMap)
	initialState.initial = true
	//pridame do Dstate

	dStates[strMapToStr(initialState.id)] = initialState

	stack := make([]dState, 0)
	stack = append(stack, initialState)
	var state dState

	for ; len(stack) > 0; {
		state, stack = stack[0], stack[1:]
		
		for _, nextStates := range state.next {
			strId := strMapToStr(nextStates)
			if _, ok := dStates[strId]; !ok {
				newDstate := getDfaStateClosure(nStates, nextStates)
				dStates[strId] = newDstate
				stack = append(stack, newDstate)
			}
		}
	}

	result := make([]dState, 0)
	for _, state := range dStates {
		result = append(result, state)
	}
	return result, alpha
}

func main() {
	states, alphabet, initial := readNStates("./in4.txt")
	dStates, alphabet := makeDFA(states, alphabet, initial)

	writeDStates(dStates, alphabet, "./out1.txt")

}
