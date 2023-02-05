package data

import (
	"campyuk-api/features/image"
	"errors"

	"gorm.io/gorm"
)

type imageData struct {
	db *gorm.DB
}

func New(db *gorm.DB) image.ImageData {
	return &imageData{db: db}
}

func (id *imageData) Add(userID uint, core image.Core) error {
	model := ToData(core)
	tx := id.db.Create(&model)
	if tx.Error != nil {
		return tx.Error
	}

	if !id.checkOwner(userID, model.ID) {
		tx.Rollback()
		return errors.New("access is denied due to invalid credential")
	}

	tx.Commit()

	return nil
}

func (id *imageData) Update(userID uint, imageID uint, core image.Core) error {
	if !id.checkOwner(userID, imageID) {
		return errors.New("access is denied due to invalid credential")
	}

	model := ToData(core)
	tx := id.db.Where("image_id = ?", imageID).Updates(&model)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (id *imageData) Delete(userID uint, imageID uint) error {
	if !id.checkOwner(userID, imageID) {
		return errors.New("access is denied due to invalid credential")
	}

	tx := id.db.Delete(&Image{}, imageID)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// ------------------------
// Functions that are not include in the contract
// ------------------------

func (id *imageData) checkOwner(hostID uint, imageID uint) bool {
	var userID uint
	query := "SELECT users.id from images JOIN camps ON camps.id = images.camp_id JOIN users ON users.id = camps.host_id WHERE images.id = ?"
	tx := id.db.Raw(query, imageID).First(&userID)
	if tx.Error != nil {
		return false
	}

	if userID != hostID {
		return false
	}

	return true
}
