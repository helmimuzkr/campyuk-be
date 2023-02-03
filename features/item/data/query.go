package data

import (
	"campyuk-api/features/item"
	"errors"
	"log"

	"gorm.io/gorm"
)

type itemData struct {
	db *gorm.DB
}

func New(db *gorm.DB) item.ItemData {
	return &itemData{
		db: db,
	}
}

func (id *itemData) checkOwner(hostID uint, itemID uint) bool {
	var userID uint
	query := "SELECT users.id from items JOIN camps ON camps.id = items.camp_id JOIN users ON users.id = camps.host_id WHERE items.id = ?"
	tx := id.db.Raw(query, itemID).First(&userID)
	if tx.Error != nil {
		return false
	}

	if userID != hostID {
		return false
	}

	return true
}

func (id *itemData) Add(userID uint, campID uint, addItem item.Core) (item.Core, error) {
	cnv := CoreToData(addItem)
	cnv.CampID = campID
	tx := id.db.Create(&cnv)
	if tx.Error != nil {
		log.Println("query error", tx.Error.Error())
		return item.Core{}, errors.New("querry error, fail to add item")
	}

	if !id.checkOwner(userID, cnv.ID) {
		tx.Rollback()
		return item.Core{}, errors.New("access is denied due to invalid credential")
	}

	tx.Commit()

	addItem.ID = cnv.ID
	return addItem, nil
}

func (id *itemData) Update(userID uint, itemID uint, updateData item.Core) (item.Core, error) {
	if !id.checkOwner(userID, itemID) {
		return item.Core{}, errors.New("access is denied due to invalid credential")
	}

	cnv := CoreToData(updateData)
	qry := id.db.Model(&Item{}).Where("id = ?", itemID).Updates(&cnv)

	affrows := qry.RowsAffected
	if affrows <= 0 {
		log.Println("no rows affected")
		return item.Core{}, errors.New("no item updated")
	}

	err := qry.Error
	if err != nil {
		log.Println("update item query error", err.Error())
		return item.Core{}, err
	}

	return ToCore(cnv), nil
}

func (id *itemData) Delete(userID uint, itemID uint) error {
	if !id.checkOwner(userID, itemID) {
		return errors.New("access is denied due to invalid credential")
	}

	data := Item{}
	qry := id.db.Where("id = ?", itemID).Delete(&data)

	affrows := qry.RowsAffected
	if affrows <= 0 {
		log.Println("no rows affected")
		return errors.New("no item deleted")
	}

	err := qry.Error
	if err != nil {
		log.Println("delete item query error", err.Error())
		return errors.New("failed to delete item")
	}

	return nil
}
