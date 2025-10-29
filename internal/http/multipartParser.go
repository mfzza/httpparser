package httpParser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func (hp *httpParser) parseMultipartBody(r *bufio.Reader) error {
	ct := strings.Split(hp.header["Content-Type"], ";")

	if len(ct) != 2 {
		return fmt.Errorf("Content-Type is multipart/form-data, but boundary not found")
	}
	boundaryStr := strings.TrimSpace(ct[1])
	if !strings.HasPrefix(boundaryStr, "boundary=") {
		return fmt.Errorf("Content-Type is multipart/form-data, but boundary is invalid: %q", boundaryStr)
	}
	boundaryStr = strings.TrimPrefix(boundaryStr, "boundary=")

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
				form, ok, err := assignMultipart(part)
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

func assignMultipart(part []byte) (multipart, bool, error) {
	form := multipart{}
	read := bufio.NewReader(bytes.NewReader(part))

	header, _, err := parseHeader(read)
	if err != nil {
		return form, false, err
	}

	// skip if doesnt have Content-Disposition field
	cd := strings.TrimSpace(header["Content-Disposition"])
	if cd == "" {
		return form, false, nil
	}

	var name, filename string
	fields := strings.Split(cd, ";")
	// if it less than 2, which is it only contain text "form-data"
	if len(fields) < 2 {
		return form, false, nil
	}
	// if the first fields is not "form-data"
	if !strings.Contains(fields[0], "form-data") {
		return form, false, nil
	}
	// if the second fields is not prefix "name="
	fields[1] = strings.TrimSpace(fields[1])
	if !strings.HasPrefix(fields[1], "name=") {
		return form, false, nil
	}

	// process name field
	name = strings.TrimPrefix(fields[1], "name=")

	// process filename field, if exists
	if len(fields) >= 3 {
		filename = strings.TrimPrefix(fields[1], "filename=")
	}

	var ct string
	if len(header["Content-Type"]) == 0 {
		ct = ""
	} else {
		ct = header["Content-Type"]
	}

	value, err := io.ReadAll(read)
	if err != nil {
		return form, false, err
	}

	form = multipart{name: name, filename: filename, contentType: ct, value: value}
	return form, true, nil
}
