package httpParser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func (hp *httpParser) parseMultipartBody(r *bufio.Reader) (error) {
	ct := strings.Split(hp.header["Content-Type"], ";")
	boundaryStr := strings.TrimPrefix(strings.TrimSpace(ct[1]), "boundary=")

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
			part := bytes.TrimPrefix(buffer[:idx], []byte("\r\n"))
			part = bytes.TrimPrefix(buffer[:idx], []byte("\n"))
			// make sure it not empty
			if len(part) != 0 {
				form, err := assignMultipart(part)
				if err != nil {
					return err
				}
				hp.forms = append(hp.forms, form)
			}
			buffer = buffer[idx+len(boundary):]

			if isBoundaryEnd(r) {
				break
			}
		}
	}

	return nil
}

func assignMultipart(part []byte) (multipart, error) {
	form := multipart{}
	read := bufio.NewReader(bytes.NewReader(part))

	header, _, err := parseHeader(read)
	if err != nil {
		return form, err
	}
	value, _ := io.ReadAll(read)

	name, filename := parseContentDisposition(header["Content-Disposition"])
	var ct string
	if len(header["Content-Type"]) == 0 {
		ct = ""
	} else {
		ct = header["Content-Type"]
	}
	form = multipart{name: name, filename: filename, contentType: ct, value: value}
	return form, nil
}

func parseContentDisposition(cd string) (string, string) {
	var name, filename string
	parts := strings.Split(cd, ";")
	name = strings.TrimPrefix(strings.TrimSpace(parts[1]), "name=")
	if len(parts) > 2 {
		filename = strings.TrimPrefix(strings.TrimSpace(parts[2]), "filename=")
	}

	return name, filename
}
