package item

import "github.com/labstack/echo/v4"

type Core struct {
	ID     uint
	Name   string
	Stock  int
	Price  int
	CampID int
}

type ItemHandler interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type ItemService interface {
	Add(token interface{}, campID uint, newItem Core) (Core, error)
	Update(token interface{}, itemID uint, updateData Core) (Core, error)
	Delete(token interface{}, itemID uint) error
}

type ItemData interface {
	Add(userID uint, newItem Core) (Core, error)
	Update(userID uint, itemID uint, updateData Core) (Core, error)
	Delete(userID uint, itemID uint) error
}
