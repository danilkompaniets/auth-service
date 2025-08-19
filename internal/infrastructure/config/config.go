package config

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	App appConfig `yaml:"application"`
}

type appConfig struct {
	GrpcAddr       string         `yaml:"grpc_addr"`
	HttpAddr       string         `yaml:"http_addr"`
	PrometheusAddr string         `yaml:"prometheus_addr"`
	Database       databaseConfig `yaml:"database"`
	Env            envConfig      `yaml:"environment"`
}
type envConfig struct {
	AccessTokenSecret  string `yaml:"accessTokenSecret"`
	RefreshTokenSecret string `yaml:"refreshTokenSecret"`
	AccessTokenTTL     string `yaml:"accessTokenTTL"`
	RefreshTokenTTL    string `yaml:"refreshTokenTTL"`
}

type databaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func MustLoad() *Config {
	err := godotenv.Load(".env")
	cfgPath := os.Getenv("CONFIG_PATH")

	if cfgPath == "" {
		panic("CONFIG_PATH environment variable not set")
	}

	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer cfgFile.Close()

	var cfg Config
	yamlParser := yaml.NewDecoder(cfgFile)
	if err := yamlParser.Decode(&cfg); err != nil {
		panic(err)
	}

	cfg.App.Env.RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
	cfg.App.Env.AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

	cfg.App.GrpcAddr = os.Getenv("GRPC_ADDR")
	cfg.App.HttpAddr = os.Getenv("HTTP_ADDR")

	cfg.App.Database.Host = os.Getenv("DB_HOST")
	cfg.App.Database.Port = os.Getenv("DB_PORT")
	cfg.App.Database.Username = os.Getenv("DB_USERNAME")
	cfg.App.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.App.Database.Database = os.Getenv("DB_NAME")

	fmt.Println(cfg)

	return &cfg
}
