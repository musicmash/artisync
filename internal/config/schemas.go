package config

import "time"

type AppConfig struct {
	Log  LogConfig  `yaml:"log"`
	DB   DBConfig   `yaml:"db"`
	HTTP HTTPConfig `yaml:"http"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type HTTPConfig struct {
	IP           string        `yaml:"ip"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DBConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	Name          string `yaml:"name"`
	Login         string `yaml:"login"`
	Password      string `yaml:"password"`
	Log           bool   `yaml:"log"`
	AutoMigrate   bool   `yaml:"auto_migrate"`
	MigrationsDir string `yaml:"migrations_dir"`
}
