package data

import (
	"campyuk-api/features/item"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name   string
	Stock  int
	Price  int
	Image  string
	CampID uint
}

func ToCore(data Item) item.Core {
	return item.Core{
		ID:     data.ID,
		CampID: int(data.CampID),
		Name:   data.Name,
		Stock:  data.Stock,
		Price:  data.Price,
		Image:  data.Image,
	}
}

func CoreToData(data item.Core) Item {
	return Item{
		Model:  gorm.Model{ID: data.ID},
		CampID: uint(data.CampID),
		Name:   data.Name,
		Stock:  data.Stock,
		Price:  data.Price,
		Image:  data.Image,
	}
}
