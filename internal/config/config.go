package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func LoadFromFile(configPath string) (*AppConfig, error) {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// replace ${ENV_NAME} in file with value from the environment
	b = []byte(os.ExpandEnv(string(b)))

	return LoadFromBytes(b)
}

func LoadFromBytes(val []byte) (*AppConfig, error) {
	config := AppConfig{}
	if err := yaml.Unmarshal(val, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *AppConfig) Dump() string {
	b, _ := yaml.Marshal(c)

	return string(b)
}

func (db *DBConfig) GetConnString() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v sslmode=disable password=%v",
		db.Host, db.Port, db.Login, db.Name, db.Password)
}
