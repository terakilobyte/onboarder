package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/terakilobyte/onboarder/globals"
)

func ParseConfigFile(configFile string) {
	// Read the config file
	if configFile == "" {
		log.Fatal("No config file specified. -c/--config is required.")
	}
	configFileContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	// Parse the config file
	err = json.Unmarshal(configFileContents, &globals.CONFIG)
	if err != nil {
		log.Fatal(err)
	}
}
