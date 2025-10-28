package main

import (
	"bufio"
	httpParser "httpparser/internal/http"
	"os"
)

func main() {
	var reader *bufio.Reader
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader = bufio.NewReader(file)
	}
	h := httpParser.NewHttpParser()
	h.Parse(reader)
	h.PrintHeader()
}
