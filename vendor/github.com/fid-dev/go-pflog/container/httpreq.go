package container

import (
	"fmt"
	"io"
)

var httpReqContainerKind = []byte("http.req")

func NewHTTPRequest(verb, uri string) *HTTPRequest {
	return &HTTPRequest{
		verb: verb,
		uri:  uri,
	}
}

type HTTPRequest struct {
	verb string
	uri  string
}

func (m *HTTPRequest) WriteTextTo(writer io.Writer) (int, error) {
	return fmt.Fprintf(writer, "%s %s", m.verb, m.uri)
}

func (m *HTTPRequest) ReadTextFrom(reader io.Reader) (int, error) {
	return 0, nil
}

func (_ HTTPRequest) Enclosed() bool {
	return true
}

func (_ HTTPRequest) Kind() []byte {
	return httpReqContainerKind
}

func (_ HTTPRequest) Multiline() bool {
	return false
}
