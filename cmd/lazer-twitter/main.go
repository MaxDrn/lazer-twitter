package main

import (
	"github.com/fid-dev/go-apiserver/pkg/util/flag"
	"github.com/fid-dev/go-apiserver/pkg/util/logs"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"lazer-twitter/pkg/cli/lazer-twitter/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())

	logs.Init()
	flag.Init()

	if err := rootCmd.Execute(); err != nil {
		glog.Fatalf("Execution failed: %s", err)
	}
}
