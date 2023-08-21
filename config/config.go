package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mallvielfrass/fmc"
)

func LookupAndParseEnvInt(envName string) (int, bool) {
	env, exists := os.LookupEnv(envName)
	if !exists {
		fmc.Printfln("#rbtError#wbt: #bbt%s", fmt.Errorf("env '%s' not found", envName).Error())
		return 0, false
	}
	parsedInt, err := strconv.Atoi(env)
	if err != nil {
		fmc.Printfln("#rbtError#wbt: #bbt%s", err.Error())
		return 0, false
	}
	return parsedInt, true
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func InitConfig(confPath string) (Config, error) {
	if confPath != "" {
		if fileExists(confPath) { //if from dot env
			if err := godotenv.Load(confPath); err != nil {
				return Config{}, fmt.Errorf("InitConfig: no '%s' file open", confPath)
			}

		} else {
			return Config{}, fmt.Errorf("InitConfig: '%s' file not exist", confPath)
		}

	}
	defaultConf := Config{
		HostPort: ":9090",
		DBHost:   "mongodb://127.0.0.1:27017/DefaultDB",
	}
	HOST_PORT, exist := os.LookupEnv("HOST_PORT")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "HOST_PORT"))
	} else {
		defaultConf.HostPort = HOST_PORT
	}
	MONGODB_HOST, exist := os.LookupEnv("MONGODB_HOST")
	if !exist {
		fmt.Printf("warn: %s\n", fmt.Errorf("env '%s' not found", "MONGODB_HOST"))
	} else {
		defaultConf.DBHost = MONGODB_HOST
	}

	return defaultConf, nil
}
