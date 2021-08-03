package service

type Config struct {
	AsanaProjectGid string
	PRLink          string
}

var conf *Config

func SetConfig(config *Config) {
	conf = config
}

func getConfig() Config {
	return *conf
}
