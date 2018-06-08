package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	config map[string]string
)

// LoadConfig load the config
func LoadConfig() string {
	result, err := ioutil.ReadFile("/etc/gowebhook/config")
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localURL, ok0 := m["requestUrl"].(string)
		localSecret, ok1 := m["secret"].(string)
		localHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		localShell, ok4 := m["script"].(string)
		if ok0 && ok1 && ok2 && ok3 && ok4 {
			config = make(map[string]string)
			config["url"] = localURL
			config["secret"] = localSecret
			config["host"] = localHost
			config["port"] = localPort
			config["shell"] = localShell
			return ""
		} else {
			return "Broken config."
		}
	} else {
		return "Can not find config file...in \"/etc/gowebhook/config\""
	}
}

// GetURL ...
func GetURL() string {
	return config["url"]
}

// GetSecret ...
func GetSecret() string {
	return config["secret"]
}

// GetHost ...
func GetHost() string {
	return config["host"]
}

// GetPort ...
func GetPort() string {
	return config["port"]
}

// GetShell ...
func GetShell() string {
	return config["shell"]
}
