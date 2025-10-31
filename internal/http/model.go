package httpParser

import (
	"bufio"
	"fmt"
)

type headerType map[string]string

type HttpParser struct {
	startLine startLineType
	header    headerType
	headerKey []string //NOTE: to know it order, maybe not really needed?
	// body      string
	body  []byte
	forms []multipart
}

type startLineType struct {
	// NOTE: assume it multipart/form-data, so it http request
	method  string
	url     string
	version string
}

// NOTE: https://datatracker.ietf.org/doc/html/rfc7578
// anatomy of multipart: each part
// - must have content-disposition header field, with type "form-data", and parameter of "name"
// - optional: parameter of "filename"
// - optional: content-type header field
// - Other header fields are generally not used and should be ignored if present
type multipart struct {
	name        string
	filename    string
	contentType string
	value       []byte
}

func NewHttpParser(r *bufio.Reader) (*HttpParser, error) {
	hp := HttpParser{header: make(headerType)}

	var err error
	err = hp.parseStartLine(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Start Line: %w", err)
	}
	hp.header, hp.headerKey, err = parseHeader(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Header field: %w", err)
	}

	err = hp.parseBody(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Body: %w", err)
	}

	return &hp, nil
}
