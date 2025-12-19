package config

import "os"

type Config struct {
	DatabaseName string
	Passwd       string
	User         string
	Host         string
}

func LoadEnv() Config {
	return Config{
		DatabaseName: os.Getenv("DATABASE_NAME"),
		Passwd:       os.Getenv("DATABASE_PASSWORD"),
		User:         os.Getenv("DATABASE_USER"),
		Host:         os.Getenv("DATABASE_HOST"),
	}
}

func SetEnvExample() {
	os.Setenv("DATABASE_NAME", "database_name")
	os.Setenv("DATABASE_PASSWORD", "oi")
	os.Setenv("DATABASE_USER", "root")
	os.Setenv("DATABASE_HOST", "127.0.0.1:3306")
}
