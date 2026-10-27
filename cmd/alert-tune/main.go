// Command alert-tune backtests Prometheus alert rules against past incidents and tunes them.
package main

import (
	"os"

	"github.com/kang0225/alert-tune/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
