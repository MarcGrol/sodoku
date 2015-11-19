package solver

import (
	"fmt"
	"os"
)

func warning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "WARNING:")
	fmt.Fprintf(os.Stderr, format, args...)
}

func debug(format string, args ...interface{}) {
	if Verbose {
		fmt.Fprintf(os.Stderr, "DEBUG:")
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
