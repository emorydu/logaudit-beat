package conf

import "sync"

var (
	conf *Config
	once sync.Once
)

type Config struct {
	LocalIP string `yaml:"local_ip"`
}

func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	// TODO:

}
