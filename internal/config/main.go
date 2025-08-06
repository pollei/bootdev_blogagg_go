package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Db_url    string
	User_name string
}

const configFileName = ".config/bootdev/gatorconfig.json"

func Read() *Config {
	var retConf Config
	confFileP, err := getConfigFilePath()
	if err != nil {
		return nil
	}
	dataBuf, err := os.ReadFile(confFileP)
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(dataBuf, &retConf); err != nil {
		return nil
	}
	return &retConf
}

func (c *Config) SetUser(s string) {
	c.User_name = s
	confFileP, err := getConfigFilePath()
	if err != nil {
		return
	}
	dataBuf, err := json.Marshal(c)
	if err != nil {
		return
	}
	os.WriteFile(confFileP, dataBuf, 0600)
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + configFileName, nil
}

// func write(cfg Config) error { return nil}
