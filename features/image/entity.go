package image

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID     uint
	CampID uint
	Image  string
}

type ImageHandler interface {
	Add() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type ImageService interface {
	Add(token interface{}, campID uint, header *multipart.FileHeader) error
	Delete(token interface{}, imageID uint) error
}

type ImageRepository interface {
	Add(userID uint, core Core) error
	Delete(usesrID uint, imageID uint) error
}

type StorageGateway interface {
	Upload(file *multipart.FileHeader) (string, error)
	Destroy(secureURL string) error
}
