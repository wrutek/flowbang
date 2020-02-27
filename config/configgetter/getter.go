package configgetter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wrutek/flowbang/config/configmodels"
	"github.com/wrutek/flowbang/settings"
	"gopkg.in/yaml.v2"
)

// GetConfiguration reads config file returns
// `Configuration` struct
func GetConfiguration() (cfg configmodels.Configuration, err error) {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	configPath := filepath.Join(currDir, settings.FileName)
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&cfg)
	if err != nil {
		return
	}
	return
}
