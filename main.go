package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"

	"github.com/lordmortis/dnsupdater/config"
)

type Options struct {
	ConfigFile string `long:"configFile" description:"path to config.yaml file" default:"config.yaml"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	options := Options{}
	parser := flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	conf, err := config.Load(&options.ConfigFile)
	if err != nil {
		fmt.Println("Unable to parse Config file")
		fmt.Println(err)
		return
	}

	fmt.Printf("Config is: %+v\n", conf)

}