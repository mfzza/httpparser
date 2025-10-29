package httpParser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

func (h *httpParser) parseMultipartBody(read *bufio.Reader) ([]multipart, error) {
	ct := h.header["Content-Type"]
	ct = strings.Split(ct[0], ";")
	boundaryStr := strings.TrimPrefix(strings.TrimSpace(ct[1]), "boundary=")

	boundary := []byte("--" + boundaryStr)
	isBoundaryEnd := func(r *bufio.Reader) bool {
		next, _ := r.Peek(2)
		return string(next) == "--"
	}

	forms := []multipart{}
	buffer := []byte{}
	var idx int
	for {
		// TODO: chunks instead of byte
		rb, err := read.ReadByte()
		if err != nil {
			return nil, err
		}

		buffer = append(buffer, rb)

		if idx = bytes.Index(buffer, boundary); idx != -1 {
			// boundary found
			// NOTE: should be both CRLF and LF
			part := bytes.TrimPrefix(buffer[:idx], []byte("\r\n"))
			part = bytes.TrimPrefix(buffer[:idx], []byte("\n"))
			// make sure it not empty
			if len(part) != 0 {
				form, err := h.assignMultipart(part)
				if err != nil {
					return nil, err
				}
				forms = append(forms, form)
			}
			buffer = buffer[idx+len(boundary):]

			if isBoundaryEnd(read) {
				break
			}
		}
	}

	return forms, nil
}

func (h *httpParser) assignMultipart(part []byte) (multipart, error) {
	form := multipart{}
	read := bufio.NewReader(bytes.NewReader(part))

	header, _, err := h.parseHeader(read)
	if err != nil {
		return form, err
	}
	value, _ := io.ReadAll(read)

	cd := header["Content-Disposition"]
	name, filename := parseContentDisposition(cd[0])
	var ct string
	if len(header["Content-Type"]) == 0 {
		ct = ""
	} else {
		ct = header["Content-Type"][0]
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
