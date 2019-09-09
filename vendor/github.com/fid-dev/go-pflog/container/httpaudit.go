package container

import (
	"fmt"
	"io"
	"time"
)

var httpAuditContainerKind = []byte("http")

func NewHTTPAudit(clientAddr, verb, uri string, status int, size int64, d time.Duration) *HTTPAudit {
	if clientAddr == "" {
		clientAddr = ""
	}

	return &HTTPAudit{
		clientAddr: clientAddr,
		verb:       verb,
		uri:        uri,
		status:     status,
		size:       size,
		duration:   d,
	}
}

type HTTPAudit struct {
	clientAddr string
	verb       string
	uri        string
	status     int
	size       int64
	duration   time.Duration
}

func (m *HTTPAudit) WriteTextTo(writer io.Writer) (int, error) {
	return fmt.Fprintf(writer, "%s - \"%s %s\" - %d %d %s", m.clientAddr, m.verb, m.uri, m.status, m.size, m.duration.String())
}

func (m *HTTPAudit) ReadTextFrom(reader io.Reader) (int, error) {
	return 0, nil
}

func (_ HTTPAudit) Enclosed() bool {
	return true
}

func (_ HTTPAudit) Kind() []byte {
	return httpAuditContainerKind
}

func (_ HTTPAudit) Multiline() bool {
	return false
}
