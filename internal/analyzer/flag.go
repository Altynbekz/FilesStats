package main

import (
	"flag"
	"fmt"
	"os"
)

type serverConfig struct {
	dir     string
	workers int
	top     int
}

func loadConfig() serverConfig {
	var config serverConfig

	flag.StringVar(&config.dir, "dir", ".", "directory to scan")
	flag.IntVar(&config.workers, "workers", 4, "number of workers")
	flag.IntVar(&config.top, "top", 10, "number of top words")

	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"Usage: %s [options]\n\nOptions:\n",
			os.Args[0],
		)

		flag.PrintDefaults()
	}

	flag.Parse()

	return config
}
