package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/terakilobyte/onboarder/globals"
)

func ParseConfigFile(configFile string) {
	// Read the config file
	configFileContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the config file
	err = json.Unmarshal(configFileContents, &globals.CONFIG)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", globals.CONFIG)
}
