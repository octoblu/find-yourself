package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/octoblu/find-yourself/scanner"
	"github.com/octoblu/find-yourself/trilateration"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("find-yourself:main")

func main() {
	app := cli.NewApp()
	app.Name = "find-yourself"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "filter, f",
			EnvVar: "FIND_YOURSELF_FILTER",
			Usage:  "Manufacturer data to filter by",
			Value:  "",
		},
	}
	app.Run(os.Args)

	select {} // Block forever
}

func run(context *cli.Context) {
	filter := getOpts(context)

	deviceScanner, err := scanner.New(filter)
	fatalIfError("scanner.New failed", err)

	swarm := trilateration.NewSwarm()
	swarm.OnLocationUpdate(func() {
		if swarm.DeviceCount() < 1 {
			return
		}
		fmt.Println("Distances: ", swarm.Distances())
	})

	deviceScanner.OnError(onError)
	deviceScanner.OnNewDeviceScanned(onNewDeviceScanned)
	deviceScanner.OnNewDeviceScanned(func(device *scanner.Device) {
		swarm.AddDevice(device)
	})
	deviceScanner.Scan()
}

func fatalIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Fatalf(msg, err.Error())
}

func getOpts(context *cli.Context) string {
	return context.String("filter")
}

func onError(err error) {
	log.Fatalln(err.Error())
}

func onNewDeviceScanned(device *scanner.Device) {
	fmt.Println("onNewDeviceScanned: ", device.String())
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
