# sodoku
Sodoku solver in go. Combines deterministic and brute force approach using recursion, go-routines and channels.
Educated guesses are run in dedicated go-routines: A bad guess will make the goroutine terminate. A good guess leading to solution will report solution back to main thread.

## Usage

    sodoku -h

or

    time sodoku < example.txt 2> /dev/null

or    

    time sodoku < hard.txt 2> /dev/null
