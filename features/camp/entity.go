package camp

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID                 uint
	VerificationStatus string
	HostName           string
	Title              string  `validate:"required,min=3"`
	Price              int     `validate:"required,number"`
	Description        string  `validate:"required,min=5"`
	Latitude           float64 `validate:"required,latitude"`
	Longitude          float64 `validate:"required,longitude"`
	Distance           int     `validate:"required,number"`
	Address            string  `validate:"required,min=5"`
	City               string  `validate:"required,min=3"`
	Document           string
	Images             []Image
	Items              []CampItem
}

type Image struct {
	ID       uint
	ImageURL string
}

type CampItem struct {
	ID        uint
	Name      string
	Stock     int
	RentPrice int
}

type CampHandler interface {
	Add() echo.HandlerFunc
	List() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	Accept() echo.HandlerFunc
	Decline() echo.HandlerFunc
}

type CampService interface {
	Add(token interface{}, newCamp Core, document *multipart.FileHeader, images []*multipart.FileHeader) error
	List(token interface{}, page int) (map[string]interface{}, []Core, error)
	GetByID(token interface{}, campID uint) (Core, error)
	Update(token interface{}, campID uint, udpateCamp Core, document *multipart.FileHeader) error
	Delete(token interface{}, campID uint) error
	RequestAdmin(token interface{}, campID uint, status string) error
}

type CampRepository interface {
	Add(userID uint, newCamp Core) error
	List(userID uint, role string, limit int, offset int) (int, []Core, error)
	GetByID(userID uint, campID uint) (Core, error)
	Update(userID uint, campID uint, updateCamp Core) error
	Delete(userID uint, campID uint) error
	RequestAdmin(campID uint, status string) error
}

type StorageGateway interface {
	Upload(file *multipart.FileHeader) (string, error)
	Destroy(secureURL string) error
}
