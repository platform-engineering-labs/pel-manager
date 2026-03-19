package main

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/platform-engineering-labs/pel-mananager/cmd/cli"
	"github.com/platform-engineering-labs/pel-mananager/cmd/ui"
)

func main() {
	if err := fang.Execute(
		context.Background(),
		cli.Root,
		fang.WithoutVersion(),
		fang.WithColorSchemeFunc(ui.FangTheme),
	); err != nil {
		os.Exit(1)
	}
}
