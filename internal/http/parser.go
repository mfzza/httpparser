package httpParser

import (
	"bufio"
	"io"
	"strings"
)

// NOTE: reused in parsing multipart
func parseHeader(read *bufio.Reader) (header, []string, error) {
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

func (h *httpParser) parseBody(read *bufio.Reader) error {
	ct := h.header["Content-Type"]
	ct = strings.Split(ct[0], ";")

	switch ct[0] {
	case "multipart/form-data":
		err := h.parseMultipartBody(read)
		if err != nil {
			return err
		}
	default:
		var err error
		h.body, err = io.ReadAll(read)
		if err != nil {
			return err
		}
	}
	return nil
}
