package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Parse() (*Config, error) {
	config := &Config{}

	configStr := viper.GetString("CONFIG")
	if configStr == "" {
		return nil, errors.New("config env is empty")
	}

	if fileExists(configStr) {
		log.Infof("loading file %s for config", configStr)

		configRaw, err := ioutil.ReadFile(configStr)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configRaw, config)

		return config, err
	}

	err := json.Unmarshal([]byte(configStr), config)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	return config, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
