package httpParser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func (hp *HttpParser) parseMultipartBody(r *bufio.Reader) error {
	boundaryStr, ok := extractBoundaryFromCt(hp.header["content-type"])
	if !ok {
		return fmt.Errorf("Content-Type is multipart/form-data, but boundary not found")
	}
	boundary := []byte("--" + boundaryStr)
	boundaryEnd := []byte("--" + boundaryStr + "--")

	buffer := make([]byte, 512)
	var temp []byte

	for {
		var stop bool
		n, err := r.Read(buffer)
		if n > 0 {
			// append previous uncomplete multipart to proceed here
			chunk := append(temp, buffer[:n]...)

			// detect boundary end
			if bytes.Contains(chunk, boundaryEnd) {
				stop = true
			}

			var parts []byte
			idx := bytes.LastIndex(chunk, boundary)
			if idx != -1 {
				// only proceed complete multipart, and store the uncomplete into temp
				parts = chunk[:idx]
				temp = chunk[idx:]
			} else {
				// if no boundary found in this iteration, empty the parts, and store all in temp
				parts = nil
				temp = chunk
			}

			if len(parts) == 0 {
				continue
			}
			for part := range bytes.SplitSeq(parts, boundary) {
				part = bytes.TrimPrefix(part, []byte("\r\n"))
				part = bytes.TrimSuffix(part, []byte("\r\n"))
				part = bytes.TrimPrefix(part, []byte("\n"))
				part = bytes.TrimSuffix(part, []byte("\n"))

				// make sure to not proceed empty part
				if len(part) == 0 {
					continue
				}

				form, ok, err := convertToMultipart(part)
				if err != nil {
					// maybe just log the error instead?
					return fmt.Errorf("Failed to convert multipart into form: %w", err)
				}
				// silently skipped if it doesnt have a MUST field
				if ok {
					hp.forms = append(hp.forms, form)
				}
			}
		}

		if stop {
			break
		}

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("Failed to read HTTP Body: %w", err)
		}
	}

	return nil
}

func convertToMultipart(part []byte) (multipart, bool, error) {
	form := multipart{}
	read := bufio.NewReader(bytes.NewReader(part))

	header, _, err := parseHeader(read)
	if err != nil {
		return form, false, fmt.Errorf("Failed to parse multipart Header field: %w", err)
	}

	cd, ok := header["content-disposition"]
	// skip if doesnt have Content-Disposition field
	if !ok {
		return form, false, nil
	}
	cd = strings.TrimSpace(cd)

	// NAME & FILENAME
	name, filename, ok := extractContentDisposition(cd)
	if !ok {
		return form, false, nil
	}

	// CONTENT-TYPE
	ct := strings.TrimSpace(header["content-type"])

	// VALUE
	value, err := io.ReadAll(read)
	if err != nil {
		return form, false, fmt.Errorf("Failed to read multipart Body: %w", err)
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

		if ok := strings.Contains(part, "form-data"); ok {
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
