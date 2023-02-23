package main

import (
	"campyuk-api/config"
	_middlewareCustom "campyuk-api/middleware"
	"campyuk-api/pkg"

	_bookingHandler "campyuk-api/features/booking/handler"
	_bookingRepo "campyuk-api/features/booking/repository"
	_bookingService "campyuk-api/features/booking/service"
	_campHandler "campyuk-api/features/camp/handler"
	_campRepo "campyuk-api/features/camp/repository"
	_campService "campyuk-api/features/camp/service"
	_imageHandler "campyuk-api/features/image/handler"
	_imageRepo "campyuk-api/features/image/repository"
	_imageService "campyuk-api/features/image/service"
	itmData "campyuk-api/features/item/data"
	itmHdl "campyuk-api/features/item/handler"
	itmSrv "campyuk-api/features/item/service"
	_userHandler "campyuk-api/features/user/handler"
	_userRepo "campyuk-api/features/user/repository"
	_userSrv "campyuk-api/features/user/service"

	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// ==================
	// Init
	// ==================
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitDB(*cfg)
	config.Migrate(db)

	// ==================
	// Init 3rd party packages
	// ==================
	v := validator.New()
	cld := pkg.NewCloudinary(cfg)
	midtransAPI := pkg.NewMidtrans(cfg)
	googleConf := pkg.NewGoogleConf(cfg)
	googleAPI := pkg.NewGoogleAPI(googleConf)

	// ==================
	// Setup services
	// ==================
	userRepo := _userRepo.New(db)
	userSrv := _userSrv.New(userRepo, v, cld, googleAPI)
	userHandler := _userHandler.New(userSrv, googleConf)

	iData := itmData.New(db)
	iSrv := itmSrv.New(iData, v)
	iHdl := itmHdl.New(iSrv)

	campRepo := _campRepo.New(db)
	campSrv := _campService.New(campRepo, v, cld)
	campHandler := _campHandler.New(campSrv)

	imageRepo := _imageRepo.New(db)
	imageSrv := _imageService.New(imageRepo, cld)
	imageHandler := _imageHandler.New(imageSrv)

	bookingRepo := _bookingRepo.New(db)
	bookingSrv := _bookingService.New(bookingRepo, midtransAPI, v, googleAPI)
	bookingHandler := _bookingHandler.New(bookingSrv)

	// ==================
	// Setup middlewares
	// ==================
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom}, method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	// ==================
	// Init routers
	// ==================
	e.POST("/register", userHandler.Register())
	e.POST("/login", userHandler.Login())
	e.GET("/auth/google", userHandler.GoogleAuth())
	e.GET("auth/google/callback", userHandler.GoogleCallback())
	e.GET("/users", userHandler.Profile(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/users", userHandler.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/users", userHandler.Delete(), middleware.JWT([]byte(config.JWT_KEY)))

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
	e.GET("/bookings/:id/reminder", bookingHandler.CreateReminder(), middleware.JWT([]byte(config.JWT_KEY)))
	e.POST("/bookings/midtrans/callback", bookingHandler.Callback())

	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}

}
