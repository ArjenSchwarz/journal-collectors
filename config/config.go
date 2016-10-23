package config

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var configFile []byte

// ParseConfig parses the config file and unmarshals it into the
// provided interface
func ParseConfig(values interface{}) error {
	filename, _ := filepath.Abs("./config.yml")
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(configFile, values); err != nil {
		return err
	}
	return nil
}
