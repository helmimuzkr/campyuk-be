package camp

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID                 uint
	VerificationStatus string
	HostName           string
	Title              string
	Price              int
	Description        string
	Latitude           float64
	Longitude          float64
	Distance           int
	Address            string
	City               string
	Document           string
	Images             []string
	Items              []CampItem
}

type CampItem struct {
	ID        uint
	Name      string
	Stock     int
	RentPrice int
	ItemImage string
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
	List(token interface{}) ([]Core, error)
	GetByID(token interface{}, campID uint) (Core, error)
	Update(token interface{}, campID uint, udpateCamp Core, document *multipart.FileHeader, images []*multipart.FileHeader) error
	Delete(token interface{}, campID uint) error
	RequestAdmin(token interface{}, campID uint) error
}

type CampData interface {
	Add(userID uint, newCamp Core) error
	List(userID uint, role string) ([]Core, error)
	GetByID(userID uint, role string, campID uint) (Core, error)
	Update(userID uint, campID uint, updateCamp Core) error
	Delete(userID uint, campID uint) error
	RequestAdmin(userID uint, campID uint) error
}
