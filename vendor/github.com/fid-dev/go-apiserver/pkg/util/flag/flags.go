package flag

import (
	"flag"

	"github.com/fid-dev/go-pflog/container"
	"github.com/fid-dev/go-pflog/log"

	"github.com/spf13/pflag"
)

func Init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	if err := flag.CommandLine.Parse([]string{}); err != nil {
		log.With(container.NewError(err, nil)).Fatal("Failed parsing arguments")
	}

	pflag.VisitAll(func(flag *pflag.Flag) {
		log.V(2).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
