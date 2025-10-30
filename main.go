package main

import (
	"bufio"
	"fmt"
	httpParser "httpparser/internal/http"
	"os"
)

func main() {
	var r *bufio.Reader
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer file.Close()

		r = bufio.NewReader(file)
	} else {
		fmt.Println("PLEASE PROVIDE THE FILE")
		fmt.Println("EXAMPLE: go run . test/testdata/multipart.txt")
		return
	}
	hp, err := httpParser.NewHttpParser(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	hp.Print()

}
