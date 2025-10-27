package main

import (
	"fmt"
	"os"
	"strings"
)

type httpParser struct {
	startLine string
	header    map[string]string
	headerKey []string // to know it order
	body      string
}

func newHttpParser() *httpParser {
	return &httpParser{header: make(map[string]string)}
}

func (h *httpParser) parseHead(head string) {
	for field := range strings.SplitSeq(head, "\n") {
		if strings.Contains(field, ":") {
			splits := strings.SplitN(field, ":", 2)
			key := strings.TrimSpace(splits[0])
			val := strings.TrimSpace(splits[1])
			h.headerKey = append(h.headerKey, key)
			h.header[key] = val
		}
	}

}

func (h *httpParser) printHead() {
	for _, k := range h.headerKey {
		fmt.Println("["+k+"]:", h.header[k])
	}
}

func main() {

	head := `POST /gnuboard4/bbs/write_update.php HTTP/1.1
Host: 192.168.100.109
Connection: keep-alive
Content-Length: 1630
Cache-Control: max-age=0
Upgrade-Insecure-Requests: 1
Origin: http://192.168.100.109`

	h := newHttpParser()
	if len(os.Args) > 1 {
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
