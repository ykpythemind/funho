package config

import "github.com/jinzhu/configor"

// Config is config
type Config struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     string `default:"3306"`
		Host     string `default:"localhost"`
	}

	Port string `required:"true"`
	Host string `required:"true"`

	ConfigFileEnv string `default:"development"`
}

// DBAddr is DBAddress
func (c Config) DBAddr() string {
	return c.DB.User + ":" + c.DB.Password + "@tcp(" + c.DB.Host + ":" + c.DB.Port + ")/" + c.DB.Name + "?charset=utf8&parseTime=True&loc=Local"
}

// Addr returns address for server
func (c Config) Addr() string {
	return c.Host + ":" + c.Port
}

// Load loads config from config.yml
func Load() Config {
	con := Config{}

	configor.Load(&con, "config.yml")
	return con
}
