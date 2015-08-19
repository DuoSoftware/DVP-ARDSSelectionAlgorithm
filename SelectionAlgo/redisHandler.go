package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var dirPath string
var redisIp string
var redisDb int
var redisPort string
var port string

const layout = "2006-01-02T15:04:05Z07:00"

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetDirPath() string {
	envPath := os.Getenv("GO_CONFIG_DIR")
	if envPath == "" {
		envPath = "./"
	}
	fmt.Println(envPath)
	return envPath
}

func GetDefaultConfig() Configuration {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("GetDefaultConfig config path: ", confPath)
	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		defconfiguration.RedisIp = "127.0.0.1"
		defconfiguration.RedisPort = "6379"
		defconfiguration.RedisDb = 6
		defconfiguration.Port = "2228"
	}

	return defconfiguration
}

func LoadDefaultConfig() {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("LoadDefaultConfig config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		redisIp = "127.0.0.1:6379"
		redisPort = "6379"
		redisDb = 6
		port = "2228"
	} else {
		redisIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
		redisPort = defconfiguration.RedisPort
		redisDb = defconfiguration.RedisDb
		port = defconfiguration.Port
	}
}

func InitiateRedis() {
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "custom-environment-variables.json")
	fmt.Println("InitiateRedis config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	envconfiguration := EnvConfiguration{}
	enverr := json.Unmarshal(content, &envconfiguration)

	if enverr != nil {
		fmt.Println("error:", enverr)
		LoadDefaultConfig()
	} else {
		var converr error
		defConfig := GetDefaultConfig()
		redisIp = os.Getenv(envconfiguration.RedisIp)
		redisPort = os.Getenv(envconfiguration.RedisPort)
		redisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
		port = os.Getenv(envconfiguration.Port)

		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisPort == "" {
			redisPort = defConfig.RedisPort
		}
		if redisDb == 0 || converr != nil {
			redisDb = defConfig.RedisDb
		}
		if port == "" {
			port = defConfig.Port
		}

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
	}

	fmt.Println("RedisIp:", redisIp)
	fmt.Println("RedisDb:", redisDb)

}
