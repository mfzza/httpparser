package httpParser

import (
	"bufio"
	"strings"
)

func (h *httpParser) parseHeader(read *bufio.Reader) (header, []string, error) {
	header := make(header)
	var headerKey []string
	for {
		line, err := read.ReadString('\n')
		if err != nil {
			return nil, nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			headerKey = append(headerKey, key)
			header[key] = append(header[key], val)
		}
	}
	return header, headerKey, nil
}

// NOTE: not used
func (h *httpParser) parseBody(read *bufio.Reader) error {
	ct := h.header["Content-Type"]
	ct = strings.Split(ct[0], ";")

	switch ct[0] {
	case "multipart/form-data":
		// boundary := strings.TrimPrefix(strings.TrimSpace(ct[1]), "boundary=")
		// h.parseMultipartBody(boundary, read)
	default:
	}
	return nil
}
