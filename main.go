package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// Known memorable primes at: 10 2446
	start := 14002
	stop := 15000
	numWorkers := 3

	workChan := make(chan int)
	primesChan := make(chan int)
	nonPrimesChan := make(chan int)
	primes := []int{}
	var wg sync.WaitGroup

	// Fire up the work producer
	go workProducer(start, stop, workChan, primesChan, nonPrimesChan, &wg)

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
				fmt.Printf("\nFound prime at %v\n", prime)
				primes = append(primes, prime)
			}
		case nonPrime, ok := <-nonPrimesChan:
			if !ok {
				nonPrimesChan = nil
			} else {
				fmt.Printf("%v ", nonPrime)
			}
		}
		if primesChan == nil && nonPrimesChan == nil {
			break
		}
	}

	fmt.Printf("\nProbable primes found at:\n")
	fmt.Println(primes)
}

func workProducer(start, stop int, workChan, primesChan, nonPrimesChan chan int, wg *sync.WaitGroup) {
	for i := start; i <= stop; i++ {
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
