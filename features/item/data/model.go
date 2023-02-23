package data

import (
	booking "campyuk-api/features/booking/repository"
	"campyuk-api/features/item"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name      string
	Stock     int
	Price     int
	CampID    uint
	RentItems []booking.RentItem `gorm:"foreignKey:ItemID"`
}

func ToCore(data Item) item.Core {
	return item.Core{
		ID:     data.ID,
		CampID: int(data.CampID),
		Name:   data.Name,
		Stock:  data.Stock,
		Price:  data.Price,
	}
}

func CoreToData(data item.Core) Item {
	return Item{
		Model:  gorm.Model{ID: data.ID},
		CampID: uint(data.CampID),
		Name:   data.Name,
		Stock:  data.Stock,
		Price:  data.Price,
	}
}
