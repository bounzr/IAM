package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type FixedConfig struct {
	Server   Server     `yaml:"server"`
	Logger   zap.Config `yaml:"logger"`
	Users    Users      `yaml:"users"`
	Clients  Clients    `yaml:"clients"`
	Groups   Groups     `yaml:"groups"`
	Sessions Sessions   `yaml:"sessions"`
	Tokens   Tokens     `yaml:"tokens"`
}

var (
	IAM        = &FixedConfig{}
	CustomConf = make(map[string]interface{})
	log        *zap.Logger
)

func Init(configFilePath string) {
	log, _ := zap.NewDevelopment()
	if len(configFilePath) == 0 {
		configFilePath = "config.yml"
	}
	ymlFile, err := os.Open(configFilePath)
	if err != nil {
		log.Fatal("Init.os.Open", zap.Error(err))
		return
	}
	defer ymlFile.Close()
	ymlData, err := ioutil.ReadAll(ymlFile)
	if err != nil {
		log.Fatal("Init.ioutil.ReadAll", zap.Error(err))
		return
	}
	err = yaml.Unmarshal([]byte(ymlData), IAM)
	if err != nil {
		log.Fatal("Init.yaml.Unmarshal", zap.Error(err))
		return
	}
	err = yaml.Unmarshal([]byte(ymlData), &CustomConf)
	if err != nil {
		log.Fatal("Init.yaml.Unmarshal", zap.Error(err))
		return
	}
}
