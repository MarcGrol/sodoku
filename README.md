# sodoku
Sodoku solver in go. Combines deterministic and brute force approach using recursion, go-routines and channels.
Educated guesses are run in dedicated go-routines: A bad guess will make the goroutine terminate. A good guess leading to solution will be reported back to main thread.

## Testing

    $ go test ./...

## Installing

    $ go install ./...
    
Expect programs "cli" and "web" to be in ${GOPATH}/bin    
    
## Usage of command-line tool

    $ cli -h

or

    $ cli < data/example.txt

or    

    $ cli < data/hardest.txt

## Usage of web-service


    $ web

and fetch the steps for a hard solution using:    

    http://localhost:3000/sodoku
