package main

import (
	"bufio"
	"fmt"
	httpParser "httpparser/internal/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "USAGE\t: go run . <path-to-file>")
		fmt.Println("EXAMPLE\t: go run . test/testdata/multipart.txt")
		os.Exit(1)
	}

	for _, filepath := range os.Args[1:] {
		// file, err := os.Open(os.Args[1])
		file, err := os.Open(filepath)
		fmt.Println("****************************************************************")
		fmt.Println("FILE:", filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening a file: %v\n", err)
			continue
		}
		defer file.Close()

		r := bufio.NewReader(file)
		hp, err := httpParser.NewHttpParser(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing HTTP Message: %v\n", err)
			continue
		}

		hp.Print()
	}

}
