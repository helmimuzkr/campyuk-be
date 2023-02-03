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
	CampID int
}

func ToCore(data Item) item.Core {
	return item.Core{
		ID:     data.ID,
		Name:   data.Name,
		Stock:  data.Stock,
		Price:  data.Price,
		CampID: data.CampID,
	}
}

func CoreToData(data item.Core) Item {
	return Item{
		Model: gorm.Model{ID: data.ID},
		Name:  data.Name,
		Stock: data.Stock,
		Price: data.Price,
	}
}
