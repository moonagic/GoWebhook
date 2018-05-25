package config

import (
	"io/ioutil"
	"encoding/json"
)

var (
	config map[string]string
)

func LoadConfig() string {
	result, err := ioutil.ReadFile("/etc/gowebhook/config")
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localUrl, ok0 := m["requestUrl"].(string)
		localSecret, ok1 := m["secret"].(string)
		localHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		localShell, ok4 := m["script"].(string)
		if ok0 && ok1 && ok2 && ok3 && ok4 {
			config = make(map[string]string)
			config["url"] = localUrl
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

func GetURL() string {
	return config["url"]
}

func GetSecret() string {
	return config["secret"]
}

func GetHost() string {
	return config["host"]
}

func GetPort() string {
	return config["port"]
}

func GetShell() string {
	return config["shell"]
}