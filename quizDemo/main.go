package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// we need some input for this demo
	csvFileName := flag.String("csv", "problemSheet.csv", "a csv file contains problems and solutions")
	timeLimit := flag.Int("limit", 30, "time limit for answering questions in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open csv: %s\n", *csvFileName))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse provided csv file")
	}

	problems := parseLines(lines)

	// initiate & start timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	//fmt.Println(problems)
	correct := 0
	for i, p := range problems {
		// if we let select case determines when to terminate current quiz, then
		// there would be a bug, user can stick around some problem
		// and let timer expired, meanwhile nothing happens until he moves to next problem
		// to solve with this, we need to use go routine to help check answer
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("Your total score is %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
				fmt.Println("Correct!")
			}
		}
	}
	// finish quiz before timer expired
	fmt.Printf("Your total score is %d out of %d.\n", correct, len(problems))
}

// This is the function convert input binary stream to problem obj array
func parseLines(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		result[i] = problem{
			q: line[0],
			a: line[1],
		}
	}

	return result
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

type problem struct {
	q string
	a string
}
