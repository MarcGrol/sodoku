# Sudoku
Sudoku solver in go. Combines deterministic and brute force approach using recursion, go-routines and channels.
Educated guesses are run in dedicated go-routines: A bad guess will make the goroutine terminate. A good guess leading to a solution will be reported back to main thread.

## Testing

    $ go test ./...

## Installing

    $ go install ./...
    
Expect programs "cliSodoku" and "webSodoku" to be in ${GOPATH}/bin    
    
## Usage of command-line tool

    $ cliSodoku -h
    	Usage of cliSodoku:
		  -help
		    	Usage information
		  -solutions int
		    	Nummber of solutions to wait for before give up (default 1)
		  -timeout int
		    	Timeout in secs before give up (default 10)
		  -verbose
		    	Verbose logging to stderr


or

    $ cliSodoku < data/example.txt

or    

    $ cliSodoku data/hardest.txt
    

## Usage of web-service

    $ webSodoku -h
	    Usage of webSodoku:
		  -help
		    	Usage information
		  -port int
		    	Port web-server is listening at (default 3000)
		  -solutions int
		    	Nummber of solutions to wait for before give up (default 1)
		  -timeout int
		    	Timeout in secs before give up (default 2)
		  -verbose
		    	Verbose logging to stderr


    $ webSodoku

and fetch the steps for a hard solution using:    

    http://localhost:3000/sodoku
