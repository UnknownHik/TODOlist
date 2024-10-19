package config

import "os"

const DateFormat = "20060102"
const LimitSearch = 20

type JWTConfig struct {
	Password string
	Secret   string
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		Password: os.Getenv("TODO_PASSWORD"),
		Secret:   os.Getenv("TODO_JWT_SECRET"),
	}
}
