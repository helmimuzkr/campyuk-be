package main

import (
	"campyuk-api/config"
	_bookingData "campyuk-api/features/booking/data"
	_bookingHandler "campyuk-api/features/booking/handler"
	_bookingService "campyuk-api/features/booking/service"
	_campData "campyuk-api/features/camp/data"
	_campHandler "campyuk-api/features/camp/handler"
	_campService "campyuk-api/features/camp/service"
	_imageData "campyuk-api/features/image/data"
	_imageHandler "campyuk-api/features/image/handler"
	_imageService "campyuk-api/features/image/service"
	itmData "campyuk-api/features/item/data"
	itmHdl "campyuk-api/features/item/handler"
	itmSrv "campyuk-api/features/item/service"
	usrData "campyuk-api/features/user/data"
	usrHdl "campyuk-api/features/user/handler"
	usrSrv "campyuk-api/features/user/services"
	"campyuk-api/helper"
	_middlewareCustom "campyuk-api/middleware"

	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitDB(*cfg)
	config.Migrate(db)

	v := validator.New()
	// cld := config.NewCloudinary(*cfg)
	coreapiMidtrans := helper.NewCoreMidtrans(cfg)

	config.Migrate(db)

	// SETUP DOMAIN
	uData := usrData.New(db)
	uSrv := usrSrv.New(uData)
	uHdl := usrHdl.New(uSrv)

	iData := itmData.New(db)
	iSrv := itmSrv.New(iData)
	iHdl := itmHdl.New(iSrv)

	campData := _campData.New(db)
	campSrv := _campService.New(campData, v)
	campHandler := _campHandler.New(campSrv)

	imageData := _imageData.New(db)
	imageSrv := _imageService.New(imageData)
	imageHandler := _imageHandler.New(imageSrv)

	bookingData := _bookingData.New(db)
	bookingSrv := _bookingService.New(bookingData, coreapiMidtrans)
	bookingHandler := _bookingHandler.New(bookingSrv)

	// MIDDLEWARE
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom}, method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	// ROUTE
	e.POST("/register", uHdl.Register())
	e.POST("/login", uHdl.Login())
	e.GET("/users", uHdl.Profile(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/users", uHdl.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/users", uHdl.Delete(), middleware.JWT([]byte(config.JWT_KEY)))

	e.POST("/camps", campHandler.Add(), middleware.JWT([]byte(config.JWT_KEY)))
	e.GET("/camps", campHandler.List(), _middlewareCustom.JWTWithConfig())
	e.GET("/camps/:id", campHandler.GetByID(), _middlewareCustom.JWTWithConfig())
	e.PUT("/camps/:id", campHandler.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/camps/:id", campHandler.Delete(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/camps/:id/accept", campHandler.Accept(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/camps/:id/decline", campHandler.Decline(), middleware.JWT([]byte(config.JWT_KEY)))

	e.POST("/images", imageHandler.Add(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/images/:id", imageHandler.Delete(), middleware.JWT([]byte(config.JWT_KEY)))

	e.POST("/items", iHdl.Add(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/items/:id", iHdl.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/items/:id", iHdl.Delete(), middleware.JWT([]byte(config.JWT_KEY)))

	e.POST("/bookings", bookingHandler.Create(), middleware.JWT([]byte(config.JWT_KEY)))
	e.GET("/bookings", bookingHandler.List(), middleware.JWT([]byte(config.JWT_KEY)))
	e.GET("/bookings/:id", bookingHandler.GetByID(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("bookings/:id/accept", bookingHandler.Accept(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("bookings/:id/cancel", bookingHandler.Cancel(), middleware.JWT([]byte(config.JWT_KEY)))
	e.POST("/bookings/callback", bookingHandler.Callback())

	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}

}
