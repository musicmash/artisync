package config

type AppConfig struct {
	Log LogConfig `yaml:"log"`
	DB  DBConfig  `yaml:"db"`
}

type LogConfig struct {
	Level string `yaml:"level"`
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
