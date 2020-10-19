package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	// assert
	assert.NoError(t, os.Setenv("DB_HOST", "artisync.db"))
	assert.NoError(t, os.Setenv("DB_PORT", "5432"))
	assert.NoError(t, os.Setenv("DB_NAME", "artisync"))
	assert.NoError(t, os.Setenv("DB_USER", "artisync"))
	assert.NoError(t, os.Setenv("DB_PASSWORD", "artisync"))

	// action
	conf, err := LoadFromFile("../../artisync.example.yml")

	// assert
	assert.NoError(t, err)

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
}
