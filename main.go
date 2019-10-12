package main

import (
	"fmt"
	"github.com/lordmortis/dnsupdater/config"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	ConfigFile string `long:"configFile" description:"path to config.yaml file" default:"config.yaml"`
}

var conf *config.Config

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	conf, err = config.Load(&options.ConfigFile)
	if err != nil {
		fmt.Println("Unable to parse Config file")
		fmt.Println(err)
		return
	}

	router := gin.Default()

	router.GET("/", updateRequest)

	err = router.Run(conf.Server.String())
	if err != nil {
		fmt.Println("Unable to start server")
		fmt.Println(err)
		return
	}

}