package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/octoblu/find-yourself/scanner"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("find-yourself:main")

func main() {
	app := cli.NewApp()
	app.Name = "find-yourself"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{}
	app.Run(os.Args)

	select {} // Block forever
}

func run(context *cli.Context) {
	getOpts(context)

	deviceScanner, err := scanner.New()
	fatalIfError("scanner.New failed", err)

	deviceScanner.OnNewDeviceScanned(onNewDeviceScanned)
	deviceScanner.Scan()
}

func fatalIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Fatalf(msg, err.Error())
}

func getOpts(context *cli.Context) {
	return
}

func onNewDeviceScanned() {
	fmt.Println("onNewDeviceScanned")
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
