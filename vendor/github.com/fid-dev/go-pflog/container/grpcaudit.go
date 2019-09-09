package container

import (
	"fmt"
	"io"
	"time"
)

var grpcAuditContainerKind = []byte("grpc")

func NewGRPCAudit(clientAddr, method, status string, size int64, d time.Duration) *GRPCAudit {
	return &GRPCAudit{
		clientAddr: clientAddr,
		method:     method,
		status:     status,
		size:       size,
		duration:   d,
	}
}

type GRPCAudit struct {
	clientAddr string
	method     string
	status     string
	size       int64
	duration   time.Duration
}

func (m *GRPCAudit) WriteTextTo(writer io.Writer) (int, error) {
	return fmt.Fprintf(writer, "%s - %s %s - %d %s", m.clientAddr, m.method, m.status, m.size, m.duration.String())
}

func (m *GRPCAudit) ReadTextFrom(reader io.Reader) (int, error) {
	return 0, nil
}

func (_ GRPCAudit) Enclosed() bool {
	return true
}

func (_ GRPCAudit) Kind() []byte {
	return grpcAuditContainerKind
}

func (_ GRPCAudit) Multiline() bool {
	return false
}
