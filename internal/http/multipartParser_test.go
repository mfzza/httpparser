package httpParser

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func Test_extractBoundaryFromCt(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ct           string
		wantBoundary string
		wantExist    bool
	}{
		{
			name:         "valid Content-Type",
			ct:           "multipart/form-data; boundary=--------------------------1234567890",
			wantBoundary: "--------------------------1234567890",
			wantExist:    true,
		},
		{
			name:         "no boundary",
			ct:           "multipart/form-data; ",
			wantBoundary: "",
			wantExist:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBoundary, gotExist := extractBoundaryFromCt(tt.ct)
			if gotBoundary != tt.wantBoundary {
				t.Errorf("extractBoundaryFromCt() = %v, want %v", gotBoundary, tt.wantBoundary)
			}
			if gotExist != tt.wantExist {
				t.Errorf("extractBoundaryFromCt() = %v, want %v", gotExist, tt.wantExist)
			}
		})
	}
}

func Test_extractContentDisposition(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cd           string
		wantName     string
		wantFilename string
		wantValid    bool
	}{
		{
			name:         "valid with filename",
			cd:           `form-data; name="bf_file[]"; filename="test.txt"`,
			wantName:     "bf_file[]",
			wantFilename: "test.txt",
			wantValid:    true,
		},
		{
			name:         "valid without name",
			cd:           `form-data; name="metadata"`,
			wantName:     "metadata",
			wantFilename: "",
			wantValid:    true,
		},
		{
			name:         "missing form-data",
			cd:           `name="bf_file[]"; filename="test.txt"`,
			wantName:     "bf_file[]",
			wantFilename: "test.txt",
			wantValid:    false,
		},
		{
			name:         "missing file field",
			cd:           `form-data; filename="test.txt"`,
			wantName:     "",
			wantFilename: "test.txt",
			wantValid:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotFilename, gotValid := extractContentDisposition(tt.cd)
			if gotName != tt.wantName {
				t.Errorf("extractContentDisposition() = %v, want %v", gotName, tt.wantName)
			}
			if gotFilename != tt.wantFilename {
				t.Errorf("extractContentDisposition() = %v, want %v", gotFilename, tt.wantFilename)
			}
			if gotValid != tt.wantValid {
				t.Errorf("extractContentDisposition() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func Test_convertToMultipart(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		part      []byte
		wantForm  multipart
		wantValid bool
		wantErr   bool
	}{
		{
			name: "valid multipart with filename",
			part: []byte(
				"Content-Disposition: form-data; name=\"file\"; filename=\"test.txt\"\r\n" +
					"Content-Type: text/plain\r\n\r\n" +
					"Hello World"),
			wantForm: multipart{
				name:        "file",
				filename:    "test.txt",
				contentType: "text/plain",
				value:       []byte("Hello World"),
			},
			wantValid: true,
			wantErr:   false,
		},
		{
			name: "valid multipart without filename",
			part: []byte(
				"Content-Disposition: form-data; name=\"field\"\r\n" +
					"Content-Type: text/plain\r\n\r\n" +
					"Some value"),
			wantForm: multipart{
				name:        "field",
				filename:    "",
				contentType: "text/plain",
				value:       []byte("Some value"),
			},
			wantValid: true,
			wantErr:   false,
		},
		{
			name: "missing Content-Disposition",
			part: []byte(
				"Content-Type: text/plain\r\n\r\n" +
					"Body content"),
			wantForm:  multipart{},
			wantValid: false,
			wantErr:   false,
		},
		{
			name: "malformed header",
			part: []byte(
				"Malformed-Header\r\n\r\nBody"),
			wantForm:  multipart{},
			wantValid: false,
			wantErr:   true, // parseHeader should fail
		},
		{
			name: "empty body",
			part: []byte(
				"Content-Disposition: form-data; name=\"field\"\r\n" +
					"Content-Type: text/plain\r\n\r\n"),
			wantForm: multipart{
				name:        "field",
				filename:    "",
				contentType: "text/plain",
				value:       []byte{},
			},
			wantValid: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotForm, gotValid, gotErr := convertToMultipart(tt.part)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("convertToMultipart() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("convertToMultipart() succeeded unexpectedly")
			}
			if gotForm.name != tt.wantForm.name || gotForm.filename != tt.wantForm.filename || gotForm.contentType != tt.wantForm.contentType || !bytes.Equal(gotForm.value, tt.wantForm.value) {
				t.Errorf("convertToMultipart() = %v, want %v", gotForm, tt.wantForm)
			}
			if gotValid != tt.wantValid {
				t.Errorf("convertToMultipart() = %v, want %v", gotValid, tt.wantValid)
			}
		})
	}
}

func Test_parseMultipartBody(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		wantForms   []multipart
		wantErr     bool
	}{
		{
			name:        "single valid part",
			contentType: "multipart/form-data; boundary=abc123",
			body: `--abc123
Content-Disposition: form-data; name="field"; filename="file.txt"
Content-Type: text/plain

Hello World
--abc123--`,
			wantForms: []multipart{
				{
					name:        "field",
					filename:    "file.txt",
					contentType: "text/plain",
					value:       []byte("Hello World"),
				},
			},
			wantErr: false,
		},
		{
			name:        "multiple parts",
			contentType: "multipart/form-data; boundary=xyz789",
			body: `--xyz789
Content-Disposition: form-data; name="name1"; filename="file1.txt"
Content-Type: text/plain

Data1
--xyz789
Content-Disposition: form-data; name="name2"; filename="file2.txt"
Content-Type: text/plain

Data2
--xyz789--`,
			wantForms: []multipart{
				{name: "name1", filename: "file1.txt", contentType: "text/plain", value: []byte("Data1")},
				{name: "name2", filename: "file2.txt", contentType: "text/plain", value: []byte("Data2")},
			},
			wantErr: false,
		},
		{
			name:        "missing boundary",
			contentType: "multipart/form-data",
			body:        "",
			wantForms:   nil,
			wantErr:     true,
		},
		{
			name:        "missing content-disposition",
			contentType: "multipart/form-data; boundary=abc123",
			body: `--abc123
Content-Type: text/plain

Hello World
--abc123--`,
			wantForms: []multipart{},
			wantErr:   false,
		},
		{
			name:        "malformed headers",
			contentType: "multipart/form-data; boundary=abc123",
			body: `--abc123
Content-Disposition

Hello World
--abc123--`,
			wantForms: nil,
			wantErr:   true,
		},
		{
			name:        "empty body",
			contentType: "multipart/form-data; boundary=abc123",
			body:        "",
			wantForms:   []multipart{},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := &HttpParser{
				header: map[string]string{"content-type": tt.contentType},
			}
			r := bufio.NewReader(strings.NewReader(tt.body))

			gotErr := hp.parseMultipartBody(r)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseMultipartBody() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("parseMultipartBody() succeeded unexpectedly")
			}
			if len(hp.forms) != len(tt.wantForms) {
				t.Fatalf("parseMultipartBody() got %d forms, want %d", len(hp.forms), len(tt.wantForms))
			}
			for i, got := range hp.forms {
				want := tt.wantForms[i]
				if got.name != want.name ||
					got.filename != want.filename ||
					got.contentType != want.contentType ||
					!bytes.Equal(got.value, want.value) {
					t.Errorf("form[%d] = %+v, want %+v", i, got, want)
				}
			}
		})
	}
}
