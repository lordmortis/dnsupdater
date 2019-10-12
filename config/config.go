package config

import (
	"io/ioutil"
	"strconv"

	"gopkg.in/yaml.v2"
)

type AWSConfig struct {
	UserID 	string 	`yaml:"id"`
	Key		string 	`yaml:"key"`
}

type ServerConfig struct {
	BindAddress	string	`yaml:"bind_address"`
	Port		int		`yaml:"port"`
}

type UserConfig struct {
	Secret	string		`yaml:"secret"`
	Name	string		`yaml:"name"`
	Hosts	[]string	`yaml:"hosts"`
}

type Config struct {
	AWS 	AWSConfig 		`yaml:"aws"`
	Server 	ServerConfig	`yaml:"server"`
	Users 	[]UserConfig	`yaml:"users"`

	UsersByHostname map[string]UserConfig
}

func defaultConfig() Config {
	var config = Config{}
	config.UsersByHostname = make(map[string]UserConfig)
	return config
}

func Load(filename *string) (*Config, error) {
	fileString, err := ioutil.ReadFile(*filename)
	if err != nil {
		return nil, err
	}

	var config = defaultConfig()

	err = yaml.Unmarshal(fileString, &config)
	if err != nil {
		return nil, err
	}


	for _, user := range config.Users {
		for _, host := range user.Hosts {
			config.UsersByHostname[host] = user
		}
	}

	return &config, nil
}

func (c ServerConfig) String() string {
	return c.BindAddress + ":" + strconv.Itoa(c.Port)
}