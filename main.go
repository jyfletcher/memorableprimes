package main

import (
	"flag"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

// Known memorable primes at: 10 2446

func main() {
	start, numWorkers := getArgs()

	workChan := make(chan int)
	primesChan := make(chan int)
	nonPrimesChan := make(chan int)
	primes := []int{}
	var wg sync.WaitGroup

	// Fire up the work producer
	go workProducer(start, workChan, primesChan, nonPrimesChan, &wg)

	// Fire up the workers
	for w := 0; w < numWorkers; w++ {
		go worker(workChan, primesChan, nonPrimesChan, &wg)
	}

	// Slurp up the results
	for {
		select {
		case prime, ok := <-primesChan:
			if !ok {
				primesChan = nil
			} else {
				fmt.Printf("\nFound probable prime at %v\n", prime)
				primes = append(primes, prime)
			}
		case nonPrime, ok := <-nonPrimesChan:
			if !ok {
				nonPrimesChan = nil
			} else {
				fmt.Printf("Found composite: %v; Probably primes at ", nonPrime)
				fmt.Println(primes)
			}
		}
		if primesChan == nil && nonPrimesChan == nil {
			break
		}
	}

	fmt.Printf("\nProbable primes found at:\n")
	fmt.Println(primes)
}

func getArgs() (int, int) {
	// Get the starting number from CLI argument, or use default
	s := flag.Int("start", 2, "What number should the test start with.")
	w := flag.Int("workers", 3, "How many concurrent worker should be started.")
	flag.Parse()
	return *s, *w
}

func workProducer(start int, workChan, primesChan, nonPrimesChan chan int, wg *sync.WaitGroup) {
	for i := start; ; i++ {
		workChan <- i
	}
	close(workChan)
	wg.Wait()
	close(primesChan)
	close(nonPrimesChan)
}

func worker(workChan, primesChan, nonPrimesChan chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	bigNum := big.NewInt(0)

	for index := range workChan {
		memorable := genMemorable(index)
		bigNum, ok := bigNum.SetString(memorable, 10)
		if !ok {
			panic("Couldn't parse number!!!")
		}
		fmt.Println("Testing: ", index)
		if bigNum.ProbablyPrime(0) {
			primesChan <- index
			continue
		} else {
			nonPrimesChan <- index
		}
	}
}

func genMemorable(limit int) string {
	var bn strings.Builder
	for i := 1; i <= limit; i++ {
		sn := strconv.Itoa(i)
		bn.WriteString(sn)
	}
	for i := limit - 1; i >= 1; i-- {
		sn := strconv.Itoa(i)
		bn.WriteString(sn)
	}
	return bn.String()
}
