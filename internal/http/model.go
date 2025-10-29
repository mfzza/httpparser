package httpParser

import (
	"bufio"
	"fmt"
)

type header map[string][]string

type httpParser struct {
	startLine string
	header    header
	headerKey []string //NOTE: to know it order, maybe not really needed?
	// body      string
	body  []byte
	forms []multipart
}

type multipart struct {
	name        string
	filename    string
	contentType string
	value       []byte
}

func NewHttpParser(r *bufio.Reader) ( *httpParser, error ) {
	hp := httpParser{header: make(map[string][]string)}

	var err error
	hp.header, hp.headerKey, err = hp.parseHeader(r)
	if err != nil {
		return nil, err
	}

	hp.forms, err = hp.parseMultipartBody(r)
	if err != nil {
		return nil, err
	}

	return &hp, nil
}

func (h *httpParser) printHeaderOrdered() {
	fmt.Println("============== HEADER ==============")
	for _, key := range h.headerKey {
		fmt.Print("- ", key, ": [", h.header[key][0], "]\n")
	}
}

func (h *httpParser) PrintHeader() {
	fmt.Println("============== HEADER ==============")
	for key, val := range h.header {
		fmt.Print("- ", key, ": [", val[0], "]\n")
	}
}

func (h *httpParser) printMultipart() {
	fmt.Println("============== HTTP DATA ==============")
	for i, form := range h.forms {
		fmt.Print("Form-Data #", i, "\n")
		fmt.Println("Field:", "["+form.name+"]")
		fmt.Println("Filename:", "["+form.filename+"]")
		fmt.Println("Content-Type:", "["+form.contentType+"]")
		fmt.Print("Value:", "["+string(form.value)+"]\n")
	}
}

func (h *httpParser) Print() {
	h.printHeaderOrdered()
	fmt.Println()
	h.printMultipart()
}
