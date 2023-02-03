package data

import (
	"campyuk-api/features/camp"

	"gorm.io/gorm"
)

type Camp struct {
	gorm.Model
	HostID             uint
	VerificationStatus string
	Title              string
	Price              int
	Description        string
	Latitude           float64
	Longtitude         float64
	Distance           int
	Address            string
	City               string
	Document           string
	CampImages         []CampImage `foreignKey:"CampID"`
}

type CampImage struct {
	gorm.Model
	CampID uint
	Image  string
}

type CampModel struct {
	ID                 uint
	VerificationStatus string
	Fullname           string
	Title              string
	Price              int
	Description        string  `gorm:"type:longtext"`
	Latitude           float64 `gorm:"type:decimal(10,15)"`
	Longitude          float64 `gorm:"type:decimal(11,14)"`
	Distance           int
	Address            string
	City               string
	Document           string
	CampImages         []string
	Items              []CampItemModel
}

type CampItemModel struct {
	ID        uint
	Name      string
	Stock     int
	RentPrice int
	ItemImage string
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
		Longtitude:         c.Longitude,
		Distance:           c.Distance,
		Address:            c.Address,
		City:               c.City,
		Document:           c.Document,
	}
}

func ToImageData(campID uint, c []string) []CampImage {
	images := []CampImage{}
	for _, v := range c {
		images = append(images, CampImage{CampID: campID, Image: v})
	}

	return images
}

func ToImageCore(ci []CampImage) []string {
	images := []string{}
	for _, v := range ci {
		images = append(images, v.Image)
	}

	return images
}

func ToItemsCore(cim []CampItemModel) []camp.CampItem {
	items := []camp.CampItem{}
	for _, v := range cim {
		i := camp.CampItem{ID: v.ID, Name: v.Name, Stock: v.Stock, RentPrice: v.RentPrice, ItemImage: v.ItemImage}
		items = append(items, i)
	}

	return items
}

func ToCampCore(cm CampModel) camp.Core {
	return camp.Core{
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
		Images:             cm.CampImages,
		Items:              ToItemsCore(cm.Items),
	}
}

func ToListCampCore(cm []CampModel) []camp.Core {
	cores := []camp.Core{}
	for _, v := range cm {
		cores = append(cores, ToCampCore(v))
	}

	return cores
}
