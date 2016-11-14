package main

import (
	"os"
	"github.com/gruntwork-io/docs/errors"
	"github.com/gruntwork-io/docs/logger"
)

// This variable is set at build time using -ldflags parameters. For more info, see:
// http://stackoverflow.com/a/11355611/483528
var VERSION string

// The main entrypoint
func main() {
	app := CreateCli(VERSION)
	err := app.Run(os.Args)

	if err != nil {
		printError(err)
		os.Exit(1)
	}
}

// Display the given error in the console
func printError(err error) {
	if os.Getenv("DOCS_PREPROCESSOR_DEBUG") != "" {
		logger.Logger.Println(errors.PrintErrorWithStackTrace(err))
	} else {
		logger.Logger.Println(err)
	}
}
