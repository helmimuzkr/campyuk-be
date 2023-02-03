package main

import (
	"campyuk-api/config"
	itmData "campyuk-api/features/item/data"
	itmHdl "campyuk-api/features/item/handler"
	itmSrv "campyuk-api/features/item/service"
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

	iData := itmData.New(db)
	iSrv := itmSrv.New(iData)
	iHdl := itmHdl.New(iSrv)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	// user
	e.POST("/register", uHdl.Register())
	e.POST("/login", uHdl.Login())
	e.GET("/users", uHdl.Profile(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/users", uHdl.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/users", uHdl.Delete(), middleware.JWT([]byte(config.JWT_KEY)))

	// camp

	// item
	e.POST("/items", iHdl.Add(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/items/:id", iHdl.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/items/:id", iHdl.Delete(), middleware.JWT([]byte(config.JWT_KEY)))

	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}

}
