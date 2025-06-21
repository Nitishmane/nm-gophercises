package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const questionTimeout = 5 * time.Second

func main() {

	csvFile, err := os.Open("input.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	// First, read all records to get the total count
	csvReader := csv.NewReader(csvFile)
	var records [][]string
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, record)
	}

	// Close and reopen the file for the quiz
	csvFile.Close()
	csvFile, err = os.Open("input.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader = csv.NewReader(csvFile)
	questionNumber := 1
	correctAnswers := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\nQuestion %d: %s\n", questionNumber, record[0])
		fmt.Printf("You have %d seconds to answer...\n", int(questionTimeout.Seconds()))

		// Create channels for communication between goroutines
		answerChan := make(chan string)
		timerChan := make(chan bool)

		// Goroutine to read user input
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answer := scanner.Text()
			answerChan <- answer
		}()

		// Goroutine for timer
		go func() {
			time.Sleep(questionTimeout)
			timerChan <- true
		}()

		// Wait for either user input or timer
		select {
		case answer := <-answerChan:
			if answer == record[1] {
				fmt.Println("Correct!")
				correctAnswers++
			} else {
				fmt.Printf("Incorrect! The correct answer was: %s\n", record[1])
				fmt.Printf("Quiz ended due to incorrect answer. You got %d questions correct out of %d\n", correctAnswers, len(records))
				return
			}
		case <-timerChan:
			fmt.Printf("\nTime's up! The correct answer was: %s\n", record[1])
			fmt.Printf("Quiz ended due to timeout. You got %d questions correct out of %d\n", correctAnswers, len(records))
			return
		}

		questionNumber++
	}

	fmt.Printf("\nCongratulations! You've completed all questions! You got %d out of %d correct!\n", correctAnswers, len(records))
}
