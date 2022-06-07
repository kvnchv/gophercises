package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// read csv
	// "flag" package parses command line args
	csvFilename := flag.String("csv", "problems.csv", "a csv file formatted 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	randomFlag := flag.Bool("random", false, "randomize question order?")
	flag.Parse()

	// flags returns a pointer that must be deref
	// 'file' is an io.reader object
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	// create csv reader object
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse provided CSV file.")
	}
	// create 1d slice of problem structs
	problems := parseLines(lines)

	// initialize timer
	// fires message over a CHANNEL C when expired
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	// <-timer.C // block and wait for message from channel

	correct := 0
	lenFull := len(problems)
	var p problem
	for len(problems) > 0 {
		p, problems = popProblem(problems, *randomFlag)

		fmt.Printf("Prompt: %s = \n", p.q)
		answerCh := make(chan string)
		go func() { // call anonymous function
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, lenFull)
			return
		case answer := <-answerCh:
			if strings.TrimSpace(answer) == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, lenFull)
}

func popProblem(problems []problem, random bool) (problem, []problem) {
	// pop problem from slice of problems, if not random just take first value
	var idx int
	var p problem
	if random {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		idx = r1.Intn(len(problems))
		p = problems[idx]
		problems = rmValueRand(problems, idx)
	} else {
		p, problems = problems[0], problems[1:]
	}
	return p, problems
}
func rmValueRand(slice []problem, i int) []problem {
	// move the last value to the current position and trim off last
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

// parse 2d slice and return 1d slice of structs
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines)) // make allocates space
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

// create a struct for modularity
type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

/* ORIGINAL SOLUTION WITHOUT RANDOMIZATION
for i, p := range problems {
	fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
	answerCh := make(chan string)
	go func() { // call anonymous function
		var answer string
		fmt.Scanf("%s\n", &answer)
		answerCh <- answer
	}()

	select {
	case <-timer.C:
		fmt.Println()
		break problemloop
	case answer := <-answerCh:
		if strings.TrimSpace(answer) == p.a {
			correct++
		}
	}
}
*/
