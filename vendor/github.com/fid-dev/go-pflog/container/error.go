package container

import (
	"bytes"
	"fmt"
	"io"
)

var errorContainerKind = []byte("error")

func NewError(err error, stack []byte) *Error {
	return &Error{
		Text:  NewText(bytes.NewBufferString(err.Error()).Bytes()),
		Stack: stack,
	}
}

type Error struct {
	*Text
	Stack []byte
}

func (m *Error) WriteTextTo(writer io.Writer) (n int, err error) {
	n, err = fmt.Fprint(writer, string(m.Text.b))
	if err != nil {
		return
	}

	if len(m.Stack) != 0 {
		var n2, n3 int

		n2, err = fmt.Fprint(writer, "\n")
		if err != nil {
			return
		}

		n += n2
		n3, err = fmt.Fprint(writer, string(m.Stack))
		if err != nil {
			return
		}

		n += n3
	}
	return
}

func (m *Error) ReadTextFrom(reader io.Reader) (int, error) {
	return 0, nil
}

func (_ Error) Enclosed() bool {
	return true
}

func (_ Error) Kind() []byte {
	return errorContainerKind
}

func (m Error) Multiline() bool {
	if len(m.Stack) != 0 {
		return true
	}
	return false
}
