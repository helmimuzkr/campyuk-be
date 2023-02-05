package data

import (
	"campyuk-api/features/camp"
	"campyuk-api/features/item/data"

	"gorm.io/gorm"
)

type Camp struct {
	gorm.Model
	HostID             uint
	VerificationStatus string
	Title              string
	Price              int
	Description        string
	Latitude           float64 `gorm:"type:decimal(16,15)"`
	Longitude          float64 `gorm:"type:decimal(17,14)"`
	Distance           int
	Address            string
	City               string
	Document           string
	Images             []Image     `gorm:"foreignKey:CampID"`
	Items              []data.Item `gorm:"foreignKey:CampID"`
}

type Image struct {
	gorm.Model
	CampID uint
	Image  string
}

type CampModel struct {
	ID                 uint
	VerificationStatus string
	HostID             uint
	Fullname           string
	Title              string
	Price              int
	Description        string
	Latitude           float64
	Longitude          float64
	Distance           int
	Address            string
	City               string
	Document           string
	Images             []Image         `gorm:"-"`
	Items              []CampItemModel `gorm:"-"`
}

type CampItemModel struct {
	ID    uint
	Name  string
	Stock int
	Price int
}

func ToData(hostID uint, c camp.Core) Camp {
	return Camp{
		Model:              gorm.Model{ID: c.ID},
		HostID:             hostID,
		VerificationStatus: c.VerificationStatus,
		Title:              c.Title,
		Price:              c.Price,
		Description:        c.Description,
		Latitude:           c.Latitude,
		Longitude:          c.Longitude,
		Distance:           c.Distance,
		Address:            c.Address,
		City:               c.City,
		Document:           c.Document,
	}
}

func ToImageData(campID uint, c []camp.Image) []Image {
	images := []Image{}
	for _, v := range c {
		images = append(images, Image{CampID: campID, Image: v.ImageURL})
	}

	return images
}

func ToImageCore(ci []Image) []camp.Image {
	images := []camp.Image{}
	for _, v := range ci {
		i := camp.Image{
			ID:       v.ID,
			ImageURL: v.Image,
		}
		images = append(images, i)
	}

	return images
}

func ToItemsCore(cim []CampItemModel) []camp.CampItem {
	items := []camp.CampItem{}
	for _, v := range cim {
		i := camp.CampItem{ID: v.ID, Name: v.Name, Stock: v.Stock, RentPrice: v.Price}
		items = append(items, i)
	}

	return items
}

func ToCampCore(cm CampModel) camp.Core {
	c := camp.Core{
		ID:                 cm.ID,
		VerificationStatus: cm.VerificationStatus,
		HostName:           cm.Fullname,
		Title:              cm.Title,
		Price:              cm.Price,
		Description:        cm.Description,
		Latitude:           cm.Latitude,
		Longitude:          cm.Longitude,
		Distance:           cm.Distance,
		Address:            cm.Address,
		City:               cm.City,
		Document:           cm.Document,
		Images:             ToImageCore(cm.Images),
		Items:              ToItemsCore(cm.Items),
	}

	return c
}

func ToListCampCore(cm []CampModel) []camp.Core {
	cores := []camp.Core{}
	for _, v := range cm {
		cores = append(cores, ToCampCore(v))
	}
	return cores
}
