package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mrcrgl/bytesf"

	"github.com/fid-dev/go-pflog/container"
	"github.com/fid-dev/go-pflog/pkg/logging"
)

var bufferPool = bytesf.NewBufferPool(128, 512)

type FlagSet interface {
	IntVar(p *int, name string, value int, usage string)
}

type FlagSetWithShorthand interface {
	IntVarP(p *int, name, shorthand string, value int, usage string)
}

func New(encoder logging.Encoder, output io.Writer) *logger {
	return &logger{
		entryPool:  logging.NewEntryPool(),
		encoder:    encoder,
		output:     output,
		containers: make([]logging.Container, 0, 5),
	}
}

var _ logging.Logger = &logger{}
var noop = new(noopLogger)

type logger struct {
	level      int
	entryPool  logging.EntryPool
	encoder    logging.Encoder
	output     io.Writer
	containers []logging.Container
}

func (l *logger) SetVerbosity(level int) int {
	previous := l.level
	l.level = level
	return previous
}

func (l logger) V(level int) logging.InfoLogger {
	if l.level >= level {
		return l
	}

	return noop
}

func (l *logger) AddFlags(fs FlagSet) {
	if f, ok := fs.(FlagSetWithShorthand); ok {
		f.IntVarP(&l.level, "verbosity", "v", l.level, "logging verbosity")
	} else {
		fs.IntVar(&l.level, "v", l.level, "logging verbosity")
	}
}

func (l logger) With(containers ...logging.Container) logging.Logger {
	// TODO alloc
	l.containers = append(l.containers, containers...)

	return l
}

func (l logger) Info(s string) {
	l.logf(logging.SeverityInfo, l.containers, s)
}

func (l logger) Infof(s string, args ...interface{}) {
	l.logf(logging.SeverityInfo, l.containers, s, args...)
}

func (l logger) Warning(s string) {
	l.logf(logging.SeverityWarning, l.containers, s)
}

func (l logger) Warningf(s string, args ...interface{}) {
	l.logf(logging.SeverityWarning, l.containers, s, args...)
}

func (l logger) Error(s string) {
	l.logf(logging.SeverityError, l.containers, s)
}

func (l logger) Errorf(s string, args ...interface{}) {
	l.logf(logging.SeverityError, l.containers, s, args...)
}

func (l logger) Fatal(s string) {
	l.logf(logging.SeverityFatal, l.containers, s)
	os.Exit(2)
}

func (l logger) Fatalf(s string, args ...interface{}) {
	l.logf(logging.SeverityFatal, l.containers, s, args...)
	os.Exit(2)
}

func (_ logger) convertToMessageContainer(format string, args ...interface{}) logging.Container {
	if format == "" {
		return nil
	}

	if len(args) == 0 {
		return container.NewMessage([]byte(format))
	}

	// TODO alloc
	b := new(bytes.Buffer)
	_, _ = fmt.Fprintf(b, format, args...)

	return container.NewMessage(b.Bytes())
}

func (l logger) logf(severity logging.Severity, containers []logging.Container, format string, args ...interface{}) {
	c := l.convertToMessageContainer(format, args...)
	if c != nil {
		containers = append(containers, c)
	}

	l.log(severity, containers)
}

func (l logger) log(severity logging.Severity, containers []logging.Container) {
	// TODO alloc
	e := l.entryPool.Allocate()

	e.Severity = severity
	e.Timestamp = time.Now()
	e.Containers = containers

	l.write(e)

	l.entryPool.Release(e)
}

func (l logger) write(entry *logging.Entry) {
	b := bufferPool.Allocate()

	// TODO alloc
	err := l.encoder.Encode(entry, b)
	if err != nil {
		fmt.Printf("Log encodimg error: %s\n", err.Error())
		bufferPool.Release(b)
		return
	}

	if _, err := b.WriteTo(l.output); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "[Logger Error] Write of log stream failed: %v\n", err)
	}

	bufferPool.Release(b)
}
