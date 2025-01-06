package config

import (
	"fmt"
	"sync"

	"github.com/Snehashish1609/couponverse-api/common"
	"github.com/joho/godotenv"
)

// couponverse Application global configuration
type Config struct {
	Port     string
	Name     string
	DBConfig DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

var once sync.Once
var config *Config

func InitConfig(name, port string) {
	once.Do(func() { // singleton

		// init DB Config
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("error loading .env file: %v\n", err)
		}

		config = &Config{
			Name: name,
			Port: port,
			DBConfig: DBConfig{
				Host:     common.GetEnvOrDie("DB_HOST"),
				Port:     common.GetEnvOrDie("DB_PORT"),
				User:     common.GetEnvOrDie("DB_USER"),
				Password: common.GetEnvOrDie("DB_PASS"),
				Name:     common.GetEnvOrDie("DB_NAME"),
			},
		}
		fmt.Println("CouponVerse global config has been set!")
	})
}

func GetConfig() *Config {
	return config
}
