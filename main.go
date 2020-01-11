package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type problem struct {
	q string
	a int
}

func main() {
	filename := flag.String("csv", "questions.csv", "Input filename")
	timeLimit := flag.String("time", "30", "Time limit (in seconds)")
	shuffle := flag.Bool("shuffle", false, "Shuffle the order of quiz questions")
	flag.Parse()

	file := openFile(filename)
	csv := openCsv(file)
	problems := getProblems(csv, *shuffle)

	scanner := bufio.NewScanner(os.Stdin)
	correct := 0

	qc := make(chan bool)
	go func() {
		for _, problem := range problems {
			fmt.Printf("Q: %s? ", problem.q)
			scanner.Scan()
			input, _ := strconv.Atoi(scanner.Text())

			if input == problem.a {
				correct++
			}
		}
		qc <- true
	}()

	timeout, _ := strconv.Atoi(*timeLimit)

	select {
	case <-qc:
		fmt.Println("Completed!")
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Println("Out of time...")
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

func getProblems(csv *csv.Reader, shuffle bool) []*problem {
	records, err := csv.ReadAll()
	bailOnError(err)

	extract := func(record []string) (string, int) {
		q, a := record[0], record[1]
		numericA, _ := strconv.Atoi(a)
		return q, numericA
	}

	idx := func() []int {
		if shuffle {
			rand.Seed(time.Now().UnixNano())
			return rand.Perm(len(records))
		}

		idx := make([]int, len(records))
		for i := 0; i < len(records); i++ {
			idx[i] = i
		}
		return idx
	}()

	problems := make([]*problem, len(records))
	for i, record := range records {
		q, a := extract(record)
		problems[idx[i]] = &problem{
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
