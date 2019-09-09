package logs

import (
	"flag"
	golog "log"

	"github.com/fid-dev/go-pflog/log"
	"github.com/fid-dev/go-pflog/pkg/logging"
	"github.com/spf13/pflag"
)

const logFlushFreqFlagName = "golog-flush-frequency"

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
}

func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(logFlushFreqFlagName))
}

type pflogWriter struct{}

func (writer pflogWriter) Write(data []byte) (n int, err error) {
	log.Info(string(data))
	return len(data), nil
}

func Init() {
	golog.SetOutput(pflogWriter{})
	golog.SetFlags(0)
}

// FlushLogs flushes logs immediately.
func FlushLogs() {

}

// NewLogger creates a new golog.Logger which sends logs to glog.Info.
func NewLogger() logging.Logger {
	return log.With()
}
