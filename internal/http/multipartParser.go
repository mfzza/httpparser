package httpParser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func (hp *httpParser) parseMultipartBody(r *bufio.Reader) error {
	boundaryStr, ok := extractBoundaryFromCt(hp.header["content-type"])
	if !ok {
		return fmt.Errorf("Content-Type is multipart/form-data, but boundary not found")
	}
	boundary := []byte("--" + boundaryStr)

	isBoundaryEnd := func(r *bufio.Reader) bool {
		next, _ := r.Peek(2)
		return string(next) == "--"
	}

	buffer := []byte{}
	var idx int
	for {
		// TODO: chunks instead of byte
		rb, err := r.ReadByte()
		if err != nil {
			return err
		}

		buffer = append(buffer, rb)

		if idx = bytes.Index(buffer, boundary); idx != -1 {
			// boundary found
			// NOTE: should be both CRLF and LF
			part := buffer[:idx]
			part = bytes.TrimPrefix(part, []byte("\r\n"))
			part = bytes.TrimPrefix(part, []byte("\n"))
			// make sure it not empty
			if len(part) != 0 {
				form, ok, err := convertToMultipart(part)
				if err != nil {
					return err
				}
				if ok {
					hp.forms = append(hp.forms, form)
				}
			}
			buffer = buffer[idx+len(boundary):]

			if isBoundaryEnd(r) {
				break
			}
		}
	}

	return nil
}

func convertToMultipart(part []byte) (multipart, bool, error) {
	form := multipart{}
	read := bufio.NewReader(bytes.NewReader(part))

	header, _, err := parseHeader(read)
	if err != nil {
		return form, false, err
	}

	// skip if doesnt have Content-Disposition field
	cd := strings.TrimSpace(header["content-disposition"])
	if cd == "" {
		return form, false, nil
	}

	name, filename, ok := extractContentDisposition(cd)
	if !ok {
		return form, false, nil
	}

	var ct string
	if len(header["content-type"]) == 0 {
		ct = ""
	} else {
		ct = header["content-type"]
	}

	value, err := io.ReadAll(read)
	if err != nil {
		return form, false, err
	}

	form = multipart{name: name, filename: filename, contentType: ct, value: value}
	return form, true, nil
}

func extractBoundaryFromCt(ct string) (string, bool) {
	for part := range strings.SplitSeq(ct, ";") {
		part = strings.TrimSpace(part)
		if after, ok := strings.CutPrefix(part, "boundary="); ok {
			after = strings.Trim(after, `"`)
			return after, true
		}
	}
	return "", false
}

func extractContentDisposition(cd string) (string, string, bool) {
	var name, filename string
	var valid bool
	for part := range strings.SplitSeq(cd, ";") {
		part = strings.TrimSpace(part)

		if ok := strings.Contains(part,"form-data"); ok {
			valid = true
		}
		if after, ok := strings.CutPrefix(part, "name="); ok {
			after = strings.Trim(after, `"`)
			name = after
		}
		if after, ok := strings.CutPrefix(part, "filename="); ok {
			after = strings.Trim(after, `"`)
			filename = after
		}
	}
	return name, filename, valid && name != ""
}
