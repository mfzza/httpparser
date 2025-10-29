package httpParser

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
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

	// NOTE: use content-length to set capacity
	cl := 0
	if v, ok := h.header["Content-Length"]; ok && len(v) > 0 {
		cl, _ = strconv.Atoi(v[0])
	}
	parts := make([][]byte, 0, cl)

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
			if len(part) != 0 {
				parts = append(parts, part)
			}
			buffer = buffer[idx+len(boundary):]

			if isBoundaryEnd(read) {
				break
			}
		}
	}

	forms := []multipart{}

	for _, part := range parts {
		read = bufio.NewReader(bytes.NewReader(part))

		header, _, err := h.parseHeader(read)
		if err != nil {
			return nil, err
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
		form := multipart{name: name, filename: filename, contentType: ct, value: value}
		forms = append(forms, form)

	}

	return forms, nil
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
