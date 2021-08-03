package github

type Config struct {
	AccessToken string
}

var conf *Config

func SetConfig(c Config) {
	conf = &c
}

func GetConfig() Config {
	if conf == nil {
		panic("asana configuration is not done yet")
	}

	return *conf
}
