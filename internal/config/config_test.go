package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestConfig_Load(t *testing.T) {
	// arrange
	assert.NoError(t, os.Setenv("DB_HOST", "artisync.db"))
	assert.NoError(t, os.Setenv("DB_PORT", "5432"))
	assert.NoError(t, os.Setenv("DB_NAME", "artisync"))
	assert.NoError(t, os.Setenv("DB_USER", "artisync"))
	assert.NoError(t, os.Setenv("DB_PASSWORD", "artisync"))
	assert.NoError(t, os.Setenv("SPOTIFY_CLIENT_ID", "2c7a0f0a-29fe-4ec4-926f-1e956297af9e"))
	assert.NoError(t, os.Setenv("SPOTIFY_CLIENT_SECRET", "75f505b3-9e40-4d55-a693-1f2388d944dd"))
	expected := AppConfig{
		HTTP: HTTPConfig{
			IP:           "0.0.0.0",
			Port:         80,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		DB: DBConfig{
			Host:                  "artisync.db",
			Port:                  5432,
			Name:                  "artisync",
			Login:                 "artisync",
			Password:              "artisync",
			AutoMigrate:           true,
			MigrationsDir:         "file:///var/artisync/migrations",
			MaxOpenConnections:    10,
			MaxIdleConnections:    10,
			MaxConnectionIdleTime: 3 * time.Minute,
			MaxConnectionLifeTime: 5 * time.Minute,
		},
		MashDB: DBConfig{
			Host:                  "musicmash.db",
			Port:                  6432,
			Name:                  "musicmash",
			Login:                 "musicmash",
			Password:              "musicmash",
			MaxOpenConnections:    15,
			MaxIdleConnections:    15,
			MaxConnectionIdleTime: 5 * time.Minute,
			MaxConnectionLifeTime: 10 * time.Minute,
		},
		Log: LogConfig{
			Level: "INFO",
		},
		Spotify: Spotify{
			AuthURL:        "https://accounts.spotify.com/authorize",
			TokenURL:       "https://accounts.spotify.com/api/token",
			ClientID:       "2c7a0f0a-29fe-4ec4-926f-1e956297af9e",
			ClientSecret:   "75f505b3-9e40-4d55-a693-1f2388d944dd",
			RedirectDomain: "https://musicmash.me",
			Scopes: []string{
				"user-top-read",
				"user-follow-read",
			},
		},
	}

	// action
	conf, err := LoadFromFile("../../artisync.example.yml")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, expected, *conf)
}
