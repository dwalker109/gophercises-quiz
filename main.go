package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type problem struct {
	q string
	a int
}

func main() {
	filename := flag.String("csv", "questions.csv", "Input filename")
	flag.Parse()

	file := openFile(filename)
	csv := openCsv(file)
	problems := getProblems(csv)

	scanner := bufio.NewScanner(os.Stdin)
	correct := 0

	for _, problem := range problems {
		fmt.Printf("Q: %s? ", problem.q)
		scanner.Scan()
		input, _ := strconv.Atoi(scanner.Text())

		if input == problem.a {
			correct++
		}
	}

	fmt.Printf("Correct: %d of %d", correct, len(problems))
}

func openFile(filename *string) *os.File {
	file, err := os.Open(*filename)
	bailOnError(err)

	return file
}

func openCsv(file *os.File) *csv.Reader {
	return csv.NewReader(file)
}

func getProblems(csv *csv.Reader) []*problem {
	records, err := csv.ReadAll()
	bailOnError(err)

	extract := func(record []string) (string, int) {
		q, a := record[0], record[1]
		numericA, _ := strconv.Atoi(a)
		return q, numericA
	}

	problems := make([]*problem, len(records))
	for i, record := range records {
		q, a := extract(record)
		problems[i] = &problem{
			q: q,
			a: a,
		}
	}

	return problems
}

func bailOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
