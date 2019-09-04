package text

import (
	"bytes"
	"io"
	"sort"

	"github.com/mrcrgl/bytesf"

	"github.com/fid-dev/go-pflog/pkg/logging"

	"github.com/mrcrgl/timef"
)

type severityStringRep []byte

var (
	infoStringRep    severityStringRep = []byte("I")
	warningStringRep severityStringRep = []byte("W")
	errorStringRep   severityStringRep = []byte("E")
	fatalStringRep   severityStringRep = []byte("F")

	newline               = "\n"
	containerIndent       = "                           "
	containerContentBegin = "{"

	spaceBuffer                        = []byte(" ")
	newlineBuffer                      = []byte(newline)
	containerContentBeginBuffer        = []byte(containerContentBegin)
	containerContentEndBuffer          = []byte("}")
	lineEndBuffer                      = []byte("\u200A\n") // unicode "HAIR SPACE" and
	lineReturnByte                     = byte('\n')
	containerIndentBuffer              = []byte(containerIndent)
	lineIndentBuffer                   = []byte("                             ")
	multilineContainerContentEndBuffer = []byte(newline + containerIndent + containerContentBegin)
)

func NewEncoder() *encoder {
	return &encoder{
		bp: bytesf.NewBufferPool(2048, 4096),
	}
}

type encoder struct {
	bp bytesf.BufferPool
}

func (e *encoder) Encode(in *logging.Entry, b *bytes.Buffer) error {
	return e.encode(in, b)
}

func Encode(in *logging.Entry) ([]byte, error) {
	bs := make([]byte, 26, 256)
	b := bytes.NewBuffer(bs)
	b.Reset()

	enc := NewEncoder()
	err := enc.Encode(in, b)

	return b.Bytes(), err
}

func (e *encoder) encode(in *logging.Entry, b *bytes.Buffer) (err error) {
	switch in.Severity {
	case logging.SeverityInfo:
		b.Write(infoStringRep)
		break
	case logging.SeverityWarning:
		b.Write(warningStringRep)
		break
	case logging.SeverityError:
		b.Write(errorStringRep)
		break
	case logging.SeverityFatal:
		b.Write(fatalStringRep)
		break
	default:
		b.WriteByte('?')
	}

	b.Write(timef.FormatRFC3339(in.Timestamp))

	if len(in.Containers) > 0 {
		sortContainers(in.Containers)

		for _, c := range in.Containers {

			if c.Multiline() {
				b.Write(newlineBuffer)
				if c.Enclosed() {
					b.Write(containerIndentBuffer)
					b.Write(c.Kind())
					b.Write(containerContentBeginBuffer)
					b.Write(newlineBuffer)
				}

				// TODO intent multiline messages
				b.Write(lineIndentBuffer)
				t := e.bp.Allocate()
				_, err := c.WriteTextTo(t)
				if err != nil {
					return err
				}

				indentMultiline(b, t)
				e.bp.Release(t)

				if c.Enclosed() {
					b.Write(multilineContainerContentEndBuffer)
				}
			} else {
				b.Write(spaceBuffer)
				if c.Enclosed() {
					b.Write(c.Kind())
					b.Write(containerContentBeginBuffer)
				}

				_, err := c.WriteTextTo(b)
				if err != nil {
					return err
				}

				if c.Enclosed() {
					b.Write(containerContentEndBuffer)
				}
			}

		}
	}

	b.Write(lineEndBuffer)

	return nil
}

func indentMultiline(w io.Writer, content *bytes.Buffer) {
	for symbol := content.Next(1); len(symbol) != 0; symbol = content.Next(1) {
		w.Write(symbol)
		if symbol[0] == lineReturnByte {
			w.Write(lineIndentBuffer)
		}
	}
}

func sortContainers(containers []logging.Container) {
	sort.Slice(containers, func(i, j int) bool {
		return !containers[i].Multiline() && containers[j].Multiline() ||
			bytes.Compare(containers[i].Kind(), containers[j].Kind()) < 0
	})
}
