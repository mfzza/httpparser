package httpParser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func (hp *httpParser) parseStartLine(r *bufio.Reader) error {
	startLine, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	startLine = strings.TrimSpace(startLine)
	parts := strings.SplitN(startLine, " ", 3)
	if len(parts) != 3 {
		return fmt.Errorf("Invalid Start Line: %q", startLine)
	}
	// NOTE: trim space for the last part because it contain whitespace
	hp.startLine = startLineType{parts[0], parts[1], parts[2]}

	return nil

}

// NOTE: reused in parsing multipart
func parseHeader(r *bufio.Reader) (headerType, []string, error) {
	header := make(headerType)
	var headerKey []string
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("Invalid Header field: %q", line)
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])
		headerKey = append(headerKey, key)
		header[key] = val
	}
	return header, headerKey, nil
}

func (hp *httpParser) parseBody(r *bufio.Reader) error {
	ct := strings.Split(hp.header["content-type"], ";")

	switch ct[0] {
	case "multipart/form-data":
		err := hp.parseMultipartBody(r)
		if err != nil {
			return err
		}
	default:
		var err error
		hp.body, err = io.ReadAll(r)
		if err != nil {
			return err
		}
	}
	return nil
}
