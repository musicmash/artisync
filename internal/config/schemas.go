package config

import (
	"time"

	"golang.org/x/oauth2"
)

type AppConfig struct {
	Log     LogConfig  `yaml:"log"`
	DB      DBConfig   `yaml:"db"`
	HTTP    HTTPConfig `yaml:"http"`
	Spotify Spotify    `yaml:"spotify"`
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
	Host                  string        `yaml:"host"`
	Port                  int           `yaml:"port"`
	Name                  string        `yaml:"name"`
	Login                 string        `yaml:"login"`
	Password              string        `yaml:"password"`
	Log                   bool          `yaml:"log"`
	AutoMigrate           bool          `yaml:"auto_migrate"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	MaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time"`
	MigrationsDir         string        `yaml:"migrations_dir"`
}

type Spotify struct {
	AuthURL      string   `yaml:"auth_url"`
	TokenURL     string   `yaml:"token_url"`
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	RedirectURL  string   `yaml:"redirect_url"`
	Scopes       []string `yaml:"scopes"`
}

func (s *Spotify) GetOAuthConfig() *oauth2.Config {
	conf := oauth2.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		RedirectURL:  s.RedirectURL,
		Scopes:       s.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  s.AuthURL,
			TokenURL: s.TokenURL,
		},
	}

	return &conf
}
