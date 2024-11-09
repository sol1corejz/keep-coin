package config

import (
	"flag"
	"os"
)

var (
	RunAddress  string
	DatabaseURI string
	LogLevel    string
)

func ParseFlags() {

	flag.StringVar(&RunAddress, "a", ":8080", "address to run server")
	flag.StringVar(&DatabaseURI, "d", "", "database uri")
	flag.StringVar(&LogLevel, "l", "info", "log level")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		RunAddress = envRunAddr
	}
	if databaseURI := os.Getenv("DATABASE_URI"); databaseURI != "" {
		DatabaseURI = databaseURI
	}

}
