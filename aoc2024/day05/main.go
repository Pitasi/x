package main

import (
	_ "embed"
	"fmt"
)

//go:embed input.txt
var input []byte

func main() {
	tokens, err := Tokenize(input)
	if err != nil {
		panic(err)
	}

	program, err := Parse(tokens)
	if err != nil {
		panic(err)
	}

	var (
		correctJobs []PrintJob
		wrongJobs   []PrintJob
	)

	for _, job := range program.PrintJobs {
		if IsSorted(program.Rules, job) {
			correctJobs = append(correctJobs, job)
		} else {
			wrongJobs = append(wrongJobs, job)
		}
	}

	star1(correctJobs)              // 4872
	star2(program.Rules, wrongJobs) // 5564
}

func star1(correctJobs []PrintJob) {
	printSum(correctJobs)
}

func star2(rules []PrecedenceRule, wrongJobs []PrintJob) {
	fixedJobs := make([]PrintJob, len(wrongJobs))
	for i, job := range wrongJobs {
		fixedJobs[i] = Sort(rules, job)
	}
	printSum(fixedJobs)
}

func printSum(jobs []PrintJob) {
	var sum int
	for _, job := range jobs {
		if len(job.Pages) == 0 {
			panic("got a PrintJob without pages")
		}
		mid := len(job.Pages) / 2
		sum += job.Pages[mid].Number
	}
	fmt.Println(sum)
}
