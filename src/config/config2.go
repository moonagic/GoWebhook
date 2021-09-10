package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"runtime"
)

var (
	Instance Config
)

func LoadConfig2() {
	configPath := ""
	switch runtime.GOOS {
	case "darwin":
		configPath = "config.yml"
		break
	case "linux":
		configPath = "/etc/gowebhook/config.yml"
		break

	}
	result, err := ioutil.ReadFile(configPath)
	if err == nil {
		Instance = Config{}
		err := yaml.Unmarshal(result, &Instance)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	} else {
		log.Fatalf("error: %v", err)
	}
}

type Config struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Servers    []Server `yaml:"servers,flow"`
}

type Server struct {
	RequestUrl string `yaml:"requestUrl"`
	Secret string `yaml:"secret"`
	Script string `yaml:"script"`
	UUID string
}