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
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type ImageService interface {
	Add(token interface{}, campID uint, header *multipart.FileHeader) error
	Update(token interface{}, imageID uint, header *multipart.FileHeader) error
	Delete(token interface{}, imageID uint) error
}

type ImageData interface {
	Add(userID uint, core Core) error
	Update(usesrID uint, imageID uint, core Core) error
	Delete(usesrID uint, imageID uint) error
}
