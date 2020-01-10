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

func main() {
	filename := flag.String("csv", "questions.csv", "Input filename")
	flag.Parse()

	file := openFile(filename)
	csv := openCsv(file)
	records := getRecords(csv)

	extract := func(record []string) (string, string) {
		return record[0], record[1]
	}

	scanner := bufio.NewScanner(os.Stdin)
	correct := 0

	for _, record := range records {
		csvQ, csvA := extract(record)

		fmt.Printf("Q: %s? ", csvQ)
		scanner.Scan()
		inputA := scanner.Text()

		intCsvA, _ := strconv.Atoi(csvA)
		intInputA, _ := strconv.Atoi(inputA)

		if intCsvA == intInputA {
			correct++
		}
	}

	fmt.Printf("Correct: %d of %d", correct, len(records))
}

struct problem {
	q: string
	a: string
}

func openFile(filename *string) *os.File {
	file, err := os.Open(*filename)
	bailOnError(err)

	return file
}

func openCsv(file *os.File) *csv.Reader {
	return csv.NewReader(file)
}

func getRecords(csv *csv.Reader) [][]string {
	records, err := csv.ReadAll()
	bailOnError(err)

	return records
}

func bailOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
