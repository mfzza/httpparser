package httpParser

import (
	"bufio"
)

type headerType map[string]string

type httpParser struct {
	startLine string
	header    headerType
	headerKey []string //NOTE: to know it order, maybe not really needed?
	// body      string
	body  []byte
	forms []multipart
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

func NewHttpParser(r *bufio.Reader) (*httpParser, error) {
	hp := httpParser{header: make(headerType)}

	var err error
	hp.header, hp.headerKey, err = parseHeader(r)
	if err != nil {
		return nil, err
	}

	err = hp.parseBody(r)
	if err != nil {
		return nil, err
	}

	return &hp, nil
}
