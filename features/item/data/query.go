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

func (id *itemData) Add(campID uint, newItem item.Core) (item.Core, error) {
	cnv := CoreToData(newItem)
	err := id.db.Create(&cnv).Error
	if err != nil {
		log.Println("query error", err.Error())
		return item.Core{}, errors.New("querry error, fail to add item")
	}

	newItem.ID = cnv.ID
	return newItem, nil
}

func (id *itemData) Update(itemID uint, campID uint, updateData item.Core) (item.Core, error) {
	cnv := CoreToData(updateData)
	qry := id.db.Model(&Item{}).Where("id = ? and camp_id = ?", itemID, campID).Updates(&cnv)

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

func (id *itemData) Delete(itemID uint, campID uint) error {
	data := Item{}
	qry := id.db.Where("id = ? and camp_id = ?", itemID, campID).Delete(&data)

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
