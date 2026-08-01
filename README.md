# Memorable Primes

---

## About

This was just some fun code testing for memorable primes as described in [The Most Wanted Prime Number - Numberphile](https://www.youtube.com/watch?v=vKlVNFOHJ9I)

From the video the only known memorable primes are at 10 and 2446.

Note that a value of 10 tests the number: 12345678910987654321

So the size of the numbers, and thus the time to test them, grows quickly.

I've tested up 15530. Nothing found so far. These 12-13k ranges with many workers can spike in memory usage, but generally stay under 100MB.

## Usage

Go 1.8+ is required. ProbablyPrime(0) before 1.8 would throw an error and 0 is chosen here to only apply the Baillie-PSW test for speed. See [math.big](https://pkg.go.dev/math/big#Int.ProbablyPrime)

The starting number to test and the number of workers can be configured as CLI arguments like this:

```sh
go run main.go --start=15530 --workers=4
```

It defaults to start=2 and workers=3 if these arguments are not specified.

## TODO

- Maybe use some even more simplified non-deterministic primality tests to more quickly get candidates that be tested more thoroughly later. I think, generally, Miller-Rabin is faster than Baille-PSW, but read more into that. The Raku version in this repo uses a Miller-Rabin test that I wrote for this type of speed increase.
- Save state so that the process can be stopped and restarted
- Related to the previous item, print the numbers that are being tested on each update, and sorted in numerical order. This allows stopping and starting manually without skipping any numbers. Some numbers can take a very long time, and some are fast, so currently it is diffifult to know if a lower number is still being tested. If you need to cancel and resume then it is not obvious where you need to resume from
- Profiling. Most of the heavy lifting is done with ProbablyPrime but the task is so CPU intensive that any savings would be valuable
