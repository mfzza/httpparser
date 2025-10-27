package main

import (
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

func (h *httpParser) parseHead(head string) {
	for field := range strings.SplitSeq(head, "\n") {
		if strings.Contains(field, ":") {
			parts := strings.SplitN(field, ":", 2)
			key := strings.TrimSpace(parts[0])
			h.headerKey = append(h.headerKey, key)

			val := strings.TrimSpace(parts[1])
			h.header[key] = append(h.header[key], val)

		}
	}

}

func (h *httpParser) printHead() {
	for _, key := range h.headerKey {
		fmt.Print(" - ", key, ": ")

		// print inline if only one
		// if len(h.header[key]) <= 1 {
		fmt.Println(h.header[key][0])

		// print multiple line if more
		// } else {
		// 	fmt.Println()
		// 	for _, val := range h.header[key] {
		// 		fmt.Println("\t-", val)
		// 	}
		// }
	}
}

func main() {

	head := `POST /gnuboard4/bbs/write_update.php HTTP/1.1
Host: 192.168.100.109
Connection: keep-alive
Content-Length: 1630
Cache-Control: max-age=0
Upgrade-Insecure-Requests: 1
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Origin: http://192.168.100.109
Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryTrctLTww4LksezWb`

	h := newHttpParser()
	if len(os.Args) > 1 {
		// stream instead turn whole of it into string?
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			panic("panik gak!")
		}
		h.parseHead(string(data))

	} else {
		h.parseHead(head)
	}

	h.printHead()
}
