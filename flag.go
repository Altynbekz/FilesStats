package main

import (
	"flag"
	"fmt"
	"os"
)

type serverConfig struct {
	port    int
	workers int
	env     string
}

func main() {
	var config serverConfig

	flag.IntVar(&config.port, "port", 8080, "port number please")
	flag.IntVar(&config.workers, "workers", 4, "number of workers please")
	flag.StringVar(&config.env, "env", "global", "please env in string")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Command line helper:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Printf("%+v\n", config)
}
