package config

type Config struct {
	Addr string
	DB   DbConfig
	Env  string
	Mail MailCfg
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
