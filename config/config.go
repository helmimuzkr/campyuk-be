package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

var (
	JWT_KEY   string
	TokenPath string = "./token.json"
)

type AppConfig struct {
	DBUser                   string
	DBPass                   string
	DBHost                   string
	DBPort                   int
	DBName                   string
	JWT_KEY                  string
	CLOUDINARY_CLOUD_NAME    string
	CLOUDINARY_API_KEY       string
	CLOUDINARY_API_SECRET    string
	MERCHANT_ID              string
	CLIENT_ID                string
	SERVER_KEY               string
	GOOGLE_REDIRECT_CALLBACK string
	GOOGLE_CLIENT_ID         string
	GOOGLE_CLIENT_SECRET     string
}

func InitConfig() *AppConfig {
	return ReadEnv()
}

func ReadEnv() *AppConfig {
	app := AppConfig{}
	isRead := true

	if val, found := os.LookupEnv("JWT_KEY"); found {
		app.JWT_KEY = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBUSER"); found {
		app.DBUser = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPASS"); found {
		app.DBPass = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBHOST"); found {
		app.DBHost = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPORT"); found {
		cnv, _ := strconv.Atoi(val)
		app.DBPort = cnv
		isRead = false
	}
	if val, found := os.LookupEnv("DBNAME"); found {
		app.DBName = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_CLOUD_NAME"); found {
		app.CLOUDINARY_CLOUD_NAME = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_API_KEY"); found {
		app.CLOUDINARY_API_KEY = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLOUDINARY_API_SECRET"); found {
		app.CLOUDINARY_API_SECRET = val
		isRead = false
	}
	if val, found := os.LookupEnv("MERCHANT_ID"); found {
		app.MERCHANT_ID = val
		isRead = false
	}
	if val, found := os.LookupEnv("CLIENT_ID"); found {
		app.CLIENT_ID = val
		isRead = false
	}
	if val, found := os.LookupEnv("SERVER_KEY"); found {
		app.SERVER_KEY = val
		isRead = false
	}
	if val, found := os.LookupEnv("GOOGLE_REDIRECT_CALLBACK"); found {
		app.GOOGLE_REDIRECT_CALLBACK = val
		isRead = false
	}
	if val, found := os.LookupEnv("GOOGLE_CLIENT_ID"); found {
		app.GOOGLE_CLIENT_ID = val
		isRead = false
	}
	if val, found := os.LookupEnv("GOOGLE_CLIENT_SECRET"); found {
		app.GOOGLE_CLIENT_SECRET = val
		isRead = false
	}

	if isRead {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error read config : ", err.Error())
			return nil
		}
		err = viper.Unmarshal(&app)
		if err != nil {
			log.Println("error parse config : ", err.Error())
			return nil
		}
	}
	JWT_KEY = app.JWT_KEY

	return &app
}
