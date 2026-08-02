package main

import (
	"flag"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type JobStatus int

const (
	StatusStarted JobStatus = iota
	StatusDone
)

type StatusUpdate struct {
	JobID  int
	Status JobStatus
}

// Known memorable primes at: 10 2446

func main() {
	start, numWorkers := getArgs()

	workChan := make(chan int)
	primesChan := make(chan int)
	nonPrimesChan := make(chan int)
	updates := make(chan StatusUpdate)
	primes := []int{}
	var wg sync.WaitGroup

	// Fire up the status monitor
	go statusMonitor(updates)

	// Fire up the work producer
	go workProducer(start, workChan, primesChan, nonPrimesChan, &wg)

	// Fire up the workers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(workChan, primesChan, nonPrimesChan, updates, &wg)
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
				fmt.Printf("Found composite: %v; Probably primes at %v\n", nonPrime, primes)
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

func worker(workChan, primesChan, nonPrimesChan chan int, updates chan<- StatusUpdate, wg *sync.WaitGroup) {
	defer wg.Done()
	bigNum := big.NewInt(0)

	for index := range workChan {
		memorable := genMemorable(index)
		bigNum, ok := bigNum.SetString(memorable, 10)
		if !ok {
			panic("Couldn't parse number!!!")
		}
		updates <- StatusUpdate{JobID: index, Status: StatusStarted}
		if bigNum.ProbablyPrime(0) {
			primesChan <- index
		} else {
			nonPrimesChan <- index
		}
		updates <- StatusUpdate{JobID: index, Status: StatusDone}
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

func statusMonitor(updates <-chan StatusUpdate) {
	// Note: time is unused for now, but let's keep track of it so we can
	// figure out a nice way to print it and give an idea of how long each
	// number takes
	active := make(map[int]time.Time) // JobID -> start time
	for u := range updates {
		switch u.Status {
		case StatusStarted:
			active[u.JobID] = time.Now()
			printStatus(active)
		case StatusDone:
			delete(active, u.JobID)
			printStatus(active)
		}
	}
}

func printStatus(active map[int]time.Time) {
	ids := make([]int, 0, len(active))
	for id := range active {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	fmt.Printf("Active: %v\n", ids)
}
