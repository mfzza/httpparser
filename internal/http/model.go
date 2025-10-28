package httpParser

import (
	"fmt"
)

type header map[string][]string

type httpParser struct {
	startLine string
	header    header
	headerKey []string //NOTE: to know it order, maybe not really needed?
	body      string
}

func NewHttpParser() *httpParser {
	return &httpParser{header: make(map[string][]string)}
}
func (h *httpParser) PrintHeaderOrdered() {
	for _, key := range h.headerKey {
		fmt.Print("- ", key, ": [", h.header[key][0], "]\n")
	}
}

func (h *httpParser) PrintHeader() {
	for key, val := range h.header {
		fmt.Print("- ", key, ": [", val[0], "]\n")
	}
}
