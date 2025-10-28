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
func (h *httpParser) PrintHeader() {
	for _, key := range h.headerKey {
		fmt.Print(" - ", key, ": ")
		fmt.Println(h.header[key][0])
	}
}
