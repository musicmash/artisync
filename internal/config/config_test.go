package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	// action
	conf, err := LoadFromFile("../../artisync.example.yml")

	// assert
	assert.NoError(t, err)

	// database section
	assert.Equal(t, conf.DB.Host, "artisync.db")
	assert.Equal(t, conf.DB.Port, 5432)
	assert.Equal(t, conf.DB.Name, "artisync+name")
	assert.Equal(t, conf.DB.Login, "artisync")
	assert.Equal(t, conf.DB.Password, "artisync")
	assert.True(t, conf.DB.AutoMigrate)
	assert.Equal(t, conf.DB.MigrationsDir, "/etc/artisync/migrations")

	// log section
	assert.Equal(t, conf.Log.Level, "INFO")
}
