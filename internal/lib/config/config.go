package config

import "time"

type Config struct {
	Addr string
	DB   DbConfig
	Env  string
	Mail MailCfg
	Auth Auth
}

type MailCfg struct {
	Addr     string
	Host     string
	Port     int
	Password string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Auth struct {
	Secret string
	Aud    string
	Iat    time.Duration
}
