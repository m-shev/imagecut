package config

import (
	"flag"
	"github.com/spf13/viper"
	"log"
)

var EnvType = struct {
	Dev  string
	QA   string
	Prod string
}{Dev: "DEV", Prod: "PROD", QA: "QA"}

var conf = &Config{isRead: false}

var configFiles = map[string]string{
	EnvType.Dev:  "dev",
	EnvType.QA:   "qa",
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
	readFlags()
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
	case EnvType.QA:
		conf.env = EnvType.QA
	default:
		conf.env = EnvType.Dev
	}
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

func readFlags() {
	size := flag.Uint("cache", 0, "cache size")
	flag.Parse()

	if *size > 0 {
		conf.CacheSize = *size
	}
}
