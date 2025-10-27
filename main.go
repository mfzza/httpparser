package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// map with key: string, value: slice of string
type Header map[string][]string

type httpParser struct {
	startLine string
	header    Header
	headerKey []string //NOTE: to know it order, maybe not really needed?
	body      string
}

func newHttpParser() *httpParser {
	return &httpParser{header: make(map[string][]string)}
}

func (h *httpParser) parseHeader(read *bufio.Reader) error {
	for {
		line, err := read.ReadString('\n')
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			h.headerKey = append(h.headerKey, key)
			h.header[key] = append(h.header[key], val)
		}
	}
	return nil
}

func (h *httpParser) printHeader() {
	for _, key := range h.headerKey {
		fmt.Print(" - ", key, ": ")
		fmt.Println(h.header[key][0])
	}
}

func main() {
	h := newHttpParser()
	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		h.parseHeader(reader)
	}
	h.printHeader()
}
