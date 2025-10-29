package httpParser

import (
	"fmt"
	"strings"
)

func (h *httpParser) Print() {
	fmt.Println("============ Start Line ============")
	h.printStartLine()
	fmt.Println("============== HEADER ==============")
	h.printHeaderOrdered()
	fmt.Println()
	fmt.Println("============== HTTP DATA ==============")
	h.printBody()
}
func (h *httpParser) printStartLine() {
	fmt.Println("- method:", h.startLine.method)
	fmt.Println("- url:", h.startLine.url)
	fmt.Println("- version:", h.startLine.version)
}

func (h *httpParser) printHeaderOrdered() {
	for _, key := range h.headerKey {
		fmt.Print("- ", key, ": [", h.header[key], "]\n")
	}
}

func (h *httpParser) printHeader() {
	for key, val := range h.header {
		fmt.Print("- ", key, ": [", val[0], "]\n")
	}
}

func (h *httpParser) printBody() {
	ct := strings.Split(h.header["Content-Type"], ";")

	switch ct[0] {
	case "multipart/form-data":
		h.printMultipart()
	default:
		fmt.Println(string(h.body))
	}

}

func (h *httpParser) printMultipart() {
	for i, form := range h.forms {
		fmt.Print("------------- Form-Data #", i+1, " -------------\n")
		fmt.Println("- Field:", "["+form.name+"]")
		fmt.Println("- Filename:", "["+form.filename+"]")
		fmt.Println("- Content-Type:", "["+form.contentType+"]")
		fmt.Println("- Value:", "["+strings.TrimSpace(string(form.value))+"]")
		// fmt.Println("Value:", form.value)
	}
}
