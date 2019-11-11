package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Http
	env    string
	isRead bool
}

type Http struct {
	Addr string
}

type Logger struct {
	Env string
}

var EnvType = struct {
	def  string
	Dev  string
	Prod string
}{Dev: "DEV", Prod: "PROD"}

var conf = &Config{env: EnvType.def, isRead: false}

var configFiles = map[string]string{
	EnvType.def:  "default",
	EnvType.Dev:  "dev",
	EnvType.Prod: "prod",
}

const (
	prefix            = "imagecut"
	envVar            = "env"
	defaultConfig     = "default"
	defaultConfigPath = "./config"
)

func GetConfig() Config {
	if !conf.isRead {
		readConfig()
	}

	return *conf
}

func GetEnv() string {
	return conf.env
}

func AddConfigPath(path string) {
	viper.AddConfigPath(path)
}

func readConfig() {
	defineEnv()
	readDefault()
	readTargetConfig()
	conf.isRead = true
}

func readDefault() {
	viper.AddConfigPath(defaultConfigPath)
	viper.SetConfigName(defaultConfig)

	read()
	unmarshal()
}

func readTargetConfig() {
	configName, ok := configFiles[conf.env]

	if ok {
		viper.SetConfigName(configName)
		read()
		unmarshal()
	} else {
		log.Fatal("Cannot read target config", configName)
	}
}

func defineEnv() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(prefix)

	env := viper.GetString(envVar)

	switch env {
	case EnvType.Prod:
		conf.env = EnvType.Prod
	case EnvType.Dev:
		conf.env = EnvType.Dev
	}
	conf.env = viper.GetString(envVar)
}

func unmarshal() {
	err := viper.Unmarshal(conf)

	if err != nil {
		log.Fatal(err)
	}
}

func read() {
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}
}
