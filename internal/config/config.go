package config

import (
	_ "github.com/mitchellh/mapstructure"
)

type Config struct {
	Port       string       `mapstructure:"port"`
	Schedulers Schedulers   `mapstructure:"schedulers"`
	TheApi     TheApiConfig `mapstructure:"the_api_config"`
}

type Schedulers struct {
	EchoJob SchedulerConfig `mapstructure:"echo_job"`
}

type SchedulerConfig struct {
	Expression string `mapstructure:"expression"`
}

type TheApiConfig struct {
	BaseUrl string `mapstructure:"base_url"`
}
