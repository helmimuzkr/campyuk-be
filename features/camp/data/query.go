package data

import (
	"campyuk-api/features/camp"

	"gorm.io/gorm"
)

type campData struct {
	db *gorm.DB
}

func New(db *gorm.DB) camp.CampData {
	return &campData{db: db}
}

func (cd *campData) Add(userID uint, newCamp camp.Core) error {
	// Create camp
	cm := ToData(userID, newCamp)
	tx := cd.db.Create(&cm)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	// Inserting image to camp
	cim := ToImageData(cm.ID, newCamp.Images)
	for _, v := range cim {
		// Kenapa pakai exec ketimbang batch create dari gorm? dikarenakan kena panic error.
		tx = tx.Exec("INSERT INTO camp_images(camp_id, image) VALUES(?, ?)", v.CampID, v.Image)
		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	tx.Commit()

	return nil
}
func (cd *campData) List(userID uint, role string) ([]camp.Core, error) {
	return nil, nil
}
func (cd *campData) GetByID(userID uint, role string, campID uint) (camp.Core, error) {
	return camp.Core{}, nil
}
func (cd *campData) Update(userID uint, campID uint, updateCamp camp.Core) error {
	return nil
}
func (cd *campData) Delete(userID uint, campID uint) error {
	return nil
}
func (cd *campData) RequestAdmin(userID uint, campID uint) error {
	return nil
}
