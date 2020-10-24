package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	// assert
	assert.NoError(t, os.Setenv("DB_HOST", "artisync.db"))
	assert.NoError(t, os.Setenv("DB_PORT", "5432"))
	assert.NoError(t, os.Setenv("DB_NAME", "artisync"))
	assert.NoError(t, os.Setenv("DB_USER", "artisync"))
	assert.NoError(t, os.Setenv("DB_PASSWORD", "artisync"))
	assert.NoError(t, os.Setenv("SPOTIFY_CLIENT_ID", "2c7a0f0a-29fe-4ec4-926f-1e956297af9e"))
	assert.NoError(t, os.Setenv("SPOTIFY_CLIENT_SECRET", "75f505b3-9e40-4d55-a693-1f2388d944dd"))

	// action
	conf, err := LoadFromFile("../../artisync.example.yml")

	// assert
	assert.NoError(t, err)

	// server section
	assert.Equal(t, "0.0.0.0", conf.HTTP.IP)
	assert.Equal(t, 80, conf.HTTP.Port)
	assert.Equal(t, 10*time.Second, conf.HTTP.ReadTimeout)
	assert.Equal(t, 10*time.Second, conf.HTTP.WriteTimeout)
	assert.Equal(t, 10*time.Second, conf.HTTP.IdleTimeout)

	// database section
	assert.Equal(t, "artisync.db", conf.DB.Host)
	assert.Equal(t, 5432, conf.DB.Port)
	assert.Equal(t, "artisync", conf.DB.Name)
	assert.Equal(t, "artisync", conf.DB.Login)
	assert.Equal(t, "artisync", conf.DB.Password)
	assert.True(t, conf.DB.AutoMigrate)
	assert.Equal(t, "file:///etc/artisync/migrations", conf.DB.MigrationsDir)

	// log section
	assert.Equal(t, conf.Log.Level, "INFO")

	// spotify section
	assert.Equal(t, "https://accounts.spotify.com/authorize", conf.Spotify.AuthURL)
	assert.Equal(t, "https://accounts.spotify.com/api/token", conf.Spotify.TokenURL)
	assert.Equal(t, "2c7a0f0a-29fe-4ec4-926f-1e956297af9e", conf.Spotify.ClientID)
	assert.Equal(t, "75f505b3-9e40-4d55-a693-1f2388d944dd", conf.Spotify.ClientSecret)
	assert.Equal(t, "https://musicmash.me/v1/callbacks/spotify/artists-sync/connect", conf.Spotify.RedirectURL)
	assert.Len(t, conf.Spotify.Scopes, 2)
}
