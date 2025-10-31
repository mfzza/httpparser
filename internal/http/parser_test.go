package httpParser

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestHttpParser_parseStartLine(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		input   string
		wantErr bool
	}{
		{
			name:    "valid start line",
			input:   "GET /index.html HTTP/1.1\r\n",
			wantErr: false,
		},
		{
			name:    "missing component (version)",
			input:   "GET /index.html\r\n",
			wantErr: true,
		},

		{
			name:    "empty line",
			input:   "\r\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := &HttpParser{}
			r := bufio.NewReader(strings.NewReader(tt.input))

			gotErr := hp.parseStartLine(r)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseStartLine() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseStartLine() succeeded unexpectedly")
			}
		})
	}
}

func Test_parseHeader(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		input   string
		want    headerType
		want2   []string
		wantErr bool
	}{
		{
			name: "valid headers",
			input: "Content-Type: text/html\r\n" +
				"User-Agent: curl/8.0\r\n" +
				"\r\n",
			want: headerType{
				"content-type": "text/html",
				"user-agent":   "curl/8.0",
			},
			want2:   []string{"content-type", "user-agent"},
			wantErr: false,
		},
		{
			name: "invalid header (missing colon)",
			input: "Content-Type text/html\r\n" +
				"\r\n",
			wantErr: true,
		},
		{
			name:    "empty Header",
			input:   "\r\n",
			want:    headerType{},
			want2:   []string{},
			wantErr: false,
		},
		{
			name:    "empty Header and Body",
			input:   "",
			want:    headerType{},
			want2:   []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.input))
			got, got2, gotErr := parseHeader(r)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseHeader() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseHeader() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeader() = got header = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("parseHeader() got headerKey = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestHttpParser_parseBody(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		ct       string
		body     string
		wantBody string
		wantErr  bool
	}{
		{
			name:     "plain text body",
			ct:       "text/plain",
			body:     "hello world",
			wantBody: "hello world",
			wantErr:  false,
		},
		{
			name:     "missing Content-Type",
			ct:       "",
			body:     "hello world",
			wantBody: "hello world",
			wantErr:  false,
		},
		{

			name:     "empty body",
			ct:       "text/plain",
			body:     "",
			wantBody: "",
			wantErr:  false,
		},
		{
			name:    "multipart body success",
			ct:      "multipart/form-data; boundary=abc",
			body:    "--abc\r\nContent-Disposition: form-data; name=\"field\"\r\n\r\nvalue\r\n--abc--",
			wantErr: false,
		},
		{
			name:    "multipart body error",
			ct:      "multipart/form-data",
			body:    "--abc\r\nsomething invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := &HttpParser{header: map[string]string{"content-type": tt.ct}}
			r := bufio.NewReader(strings.NewReader(tt.body))
			gotErr := hp.parseBody(r)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseBody() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseBody() succeeded unexpectedly")
			}
			if string(hp.body) != tt.wantBody {
				t.Errorf("form = %+v, want %+v", string(hp.body), tt.body)
			}
		})
	}
}
