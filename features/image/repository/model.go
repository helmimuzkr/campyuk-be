package repository

import (
	"campyuk-api/features/image"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	CampID uint
	Image  string
}

func ToData(data image.Core) Image {
	return Image{
		Model:  gorm.Model{ID: data.ID},
		CampID: uint(data.CampID),
		Image:  data.Image,
	}
}

func ToCore(data Image) image.Core {
	return image.Core{
		ID:     data.ID,
		CampID: uint(data.CampID),
		Image:  data.Image,
	}
}
