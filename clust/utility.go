package clust

import (
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Servers map[string]string	
}

func LoadConfig(configFile string,c *Config) error {
	fileBytes,_ := ioutil.ReadFile(configFile)
	err := json.Unmarshal(fileBytes,c)
	return err
}


