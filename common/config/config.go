package config

type Config struct {
	DefaultSchema string
	DefaultDB     string
}

var DBConfig = &Config{}
