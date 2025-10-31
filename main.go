package main

import (
	"bufio"
	"fmt"
	httpParser "httpparser/internal/http"
	"os"
)

func main() {
	warning := func() string {
		return "USAGE\t: go run . <path-to-file>\nEXAMPLE\t: go run . test/testdata/multipart.txt"
	}
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, warning())
		os.Exit(2)
	}

	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Too many arguments\n" + warning())
		os.Exit(2)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening a file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	hp, err := httpParser.NewHttpParser(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing HTTP Message: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("FILE:", os.Args[1] )
	hp.Print()
}
