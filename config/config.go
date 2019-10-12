package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AWSConfig struct {
	UserID 	string 	`yaml:"id"`
	Key		string 	`yaml:"key"`
}

type UserConfig struct {
	Secret	string		`yaml:"secret"`
	Name	string		`yaml:"name"`
	Hosts	[]string	`yaml:"hosts"`
}

type Config struct {
	AWS AWSConfig 		`yaml:"aws"`
	Users []UserConfig	`yaml:"users"`
}

func defaultConfig() Config {
	var config = Config{}
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

	return &config, nil
}