package main

import (
	"campyuk-api/config"
	usrData "campyuk-api/features/user/data"
	usrHdl "campyuk-api/features/user/handler"
	usrSrv "campyuk-api/features/user/services"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitDB(*cfg)

	// v := validator.New()
	// cld := config.NewCloudinary(*cfg)

	config.Migrate(db)

	uData := usrData.New(db)
	uSrv := usrSrv.New(uData)
	uHdl := usrHdl.New(uSrv)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	// user
	e.POST("/register", uHdl.Register())
	e.POST("/login", uHdl.Login())
	e.PUT("/users", uHdl.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/users", uHdl.Delete(), middleware.JWT([]byte(config.JWT_KEY)))
	e.GET("/users", uHdl.Profile(), middleware.JWT([]byte(config.JWT_KEY)))

	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}

}
