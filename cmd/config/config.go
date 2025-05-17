package config

import (
	"flag"
	"os"
)

var (
	RunAddress  string
	DatabaseURI string
)

func ParseFlags() {

	flag.StringVar(&RunAddress, "a", ":8080", "address to run server")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		RunAddress = envRunAddr
	}

}
