package config

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type AppConfig struct {
	Log     LogConfig  `yaml:"log"`
	DB      DBConfig   `yaml:"db"`
	MashDB  DBConfig   `yaml:"musicmash_db"`
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
	AutoMigrate           bool          `yaml:"auto_migrate"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	MaxConnectionLifeTime time.Duration `yaml:"max_connection_life_time"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time"`
	MigrationsDir         string        `yaml:"migrations_dir"`
}

type Spotify struct {
	AuthURL        string   `yaml:"auth_url"`
	TokenURL       string   `yaml:"token_url"`
	ClientID       string   `yaml:"client_id"`
	ClientSecret   string   `yaml:"client_secret"`
	RedirectDomain string   `yaml:"redirect_domain"`
	Scopes         []string `yaml:"scopes"`
}

// Returns oAuth credentials
//
// Note: RedirectURL will be empty.
// If you need to get credentials with filled RedirectURL
// just call GetOnceSyncOAuthConfig or GetDailySyncOAuthConfig methods.
func (s *Spotify) GetOAuthConfig() *oauth2.Config {
	conf := oauth2.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		Scopes:       s.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  s.AuthURL,
			TokenURL: s.TokenURL,
		},
	}

	return &conf
}

func (s *Spotify) GetOnceSyncOAuthConfig() *oauth2.Config {
	conf := s.GetOAuthConfig()

	conf.RedirectURL = fmt.Sprintf("%s/v1/artists/sync/once/connect", s.RedirectDomain)

	return conf
}

func (s *Spotify) GetDailySyncOAuthConfig() *oauth2.Config {
	conf := s.GetOAuthConfig()

	conf.RedirectURL = fmt.Sprintf("%s/v1/artists/sync/daily/connect", s.RedirectDomain)

	return conf
}
