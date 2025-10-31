package httpParser

import (
	"fmt"
	"strings"
)

func (hp *HttpParser) Print() {
	fmt.Println("============ Start Line ============")
	hp.printStartLine()
	fmt.Println("============== HEADER ==============")
	hp.printHeaderOrdered()
	fmt.Println("============== HTTP DATA ==============")
	hp.printBody()
}
func (hp *HttpParser) printStartLine() {
	fmt.Println("- method:", hp.startLine.method)
	fmt.Println("- url:", hp.startLine.url)
	fmt.Println("- version:", hp.startLine.version)
}

func (hp *HttpParser) printHeaderOrdered() {
	for _, key := range hp.headerKey {
		fmt.Print("- ", key, ": [", hp.header[key], "]\n")
	}
}

func (hp *HttpParser) printHeader() {
	for key, val := range hp.header {
		fmt.Print("- ", key, ": [", val, "]\n")
	}
}

func (hp *HttpParser) printBody() {
	ct := strings.Split(hp.header["content-type"], ";")

	switch ct[0] {
	case "multipart/form-data":
		hp.printMultipart()
	default:
		fmt.Println(string(hp.body))
	}

}

func (hp *HttpParser) printMultipart() {
	for i, form := range hp.forms {
		fmt.Print("------------- Form-Data #", i+1, " -------------\n")
		fmt.Println("- Field:", "["+form.name+"]")
		fmt.Println("- Filename:", "["+form.filename+"]")
		fmt.Println("- Content-Type:", "["+form.contentType+"]")
		fmt.Println("- Value:", "["+strings.TrimSpace(string(form.value))+"]")
		// fmt.Println("Value:", form.value)
	}
}
