package config

import (
	"io/ioutil"
	"runtime"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

var (
	Instance Config
)

func LoadConfig() {
	configPath := ""
	switch runtime.GOOS {
	case "darwin":
		configPath = "config.example.yml"
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
			color.Red("error: %v", err)
		}

		identifiers := make([]string, 0)
		secrets := make([]string, 0)
		for _, server := range Instance.Servers {
			if !contains(identifiers, server.Identifier) {
				identifiers = append(identifiers, server.Identifier)
			} else {
				color.Red("Identifiers in config have to different")
			}
			if !contains(secrets, server.Secret) {
				secrets = append(secrets, server.Secret)
			} else {
				color.Red("Secrets in config have to different")
			}
		}
	} else {
		color.Red("error: %v", err)
	}
}

func contains(array []string, target string) bool {
	for _, value := range array {
		if value == target {
			return true
		}
	}
	return false
}

type Config struct {
	Host    string   `yaml:"host"`
	Port    string   `yaml:"port"`
	Servers []Server `yaml:"servers,flow"`
}

type Server struct {
	Identifier string `yaml:"identifier"`
	RequestUrl string `yaml:"requestUrl"`
	Secret     string `yaml:"secret"`
	Script     string `yaml:"script"`
	UUID       string
}
