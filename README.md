# Memorable Primes

---

## About

This was just some fun code testing for memorable primes as described in [The Most Wanted Prime Number - Numberphile](https://www.youtube.com/watch?v=vKlVNFOHJ9I)

From the video the only known memorable primes are at 10 and 2446.

Note that a value of 10 tests the number: 12345678910987654321

So the size of the numbers, and thus the time to test them, grows quickly.

I've tested up to just under 13000 and ranges of 1000 (13000-14000) with 4 workers can take weeks. The 12-13k range with 4 works can spike in memory usage but typically hangs around 80MB.

## Usage

Go 1.8+ is required. ProbablyPrime(0) before 1.8 would throw an error and 0 is chosen here to only apply the Baillie-PSW test for speed. See [math.big](https://pkg.go.dev/math/big#Int.ProbablyPrime)

Set the start and stop values in the main function.

Adjust numWorkers to set the number of goroutines used in testing numbers in the range.

## TODO

- Maybe use some even more simplified non-deterministic primality tests to more quickly get candidates that be tested more thouroughly later
- Use command line args instead of changing the source for the start, stop, and numWorkers values
- Save state so that the process can be stopped and restarted
