package config

import (
	"flag"
	"os"
)

func getConfigFilePath() string {
	var configFile string
	flag.StringVar(&configFile, "c", "", "CONFIG")
	if len(configFile) == 0 {
		flag.StringVar(&configFile, "config", "", "CONFIG")
	}

	tmpVal := os.Getenv("CONFIG")
	if len(tmpVal) != 0 {
		configFile = tmpVal
	}
	return configFile
}
