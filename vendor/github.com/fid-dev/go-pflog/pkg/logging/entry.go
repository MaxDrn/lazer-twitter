package logging

import (
	"bytes"
	"time"
)

type Encoder interface {
	Encode(in *Entry, buffer *bytes.Buffer) error
}

type Decoder interface {
	Decode(in []byte, out *Entry) (err error)
}

type Entry struct {
	Severity   Severity
	Timestamp  time.Time
	Containers []Container
}
