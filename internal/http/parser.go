package httpParser

import (
	"bufio"
	"fmt"
	"strings"
)

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


func (h *httpParser) parseBody(read *bufio.Reader) error {
	ct := h.header["Content-Type"]
	ct = strings.SplitN(ct[0], ";", 2)

	switch ct[0] {
	case "multipart/form-data":
		// NOTE: https://datatracker.ietf.org/doc/html/rfc7578
		// anatomy of multipart: each part
		// - must have content-disposition header field, with type "form-data", and parameter of "name"
		// - optional: parameter of "filename"
		// - optional: content-type header field
		// - Other header fields are generally not used and should be ignored if present

		boundary := ct[1]
		boundary = "--" + strings.TrimPrefix(boundary, " boundary=")
		// TODO: process multipart/form-data
		for {
			line, err := read.ReadString('\n')
			if err != nil {
				return err
			}
			line = strings.TrimSpace(line)
			if line == boundary {
				continue
			} else {
				if strings.Contains(line, ":") {
					parts := strings.SplitN(line, ":", 2)
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					fmt.Println("-", key+ ":", "[" + val + "]")
				}
			}
		}
	default:
	}
	return nil
}

func (h *httpParser) Parse(read *bufio.Reader) error {
	// TODO: bundle error
	h.parseHeader(read)
	h.parseBody(read)
	return nil
}

