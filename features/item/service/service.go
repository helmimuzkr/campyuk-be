package service

import (
	"campyuk-api/features/item"
	"campyuk-api/helper"
	"errors"
	"strings"
)

type itemSrv struct {
	qry item.ItemData
}

func New(id item.ItemData) item.ItemService {
	return &itemSrv{
		qry: id,
	}
}

func (is *itemSrv) Add(token interface{}, campID uint, newItem item.Core) (item.Core, error) {
	userID, _ := helper.ExtractToken(token)
	if userID <= 0 {
		return item.Core{}, errors.New("data not found")
	}

	res, err := is.qry.Add(userID, campID, newItem)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "internal server error"
		}
		return item.Core{}, errors.New(msg)
	}

	return res, nil
}

func (is *itemSrv) Update(token interface{}, itemID uint, updateData item.Core) (item.Core, error) {
	userID, _ := helper.ExtractToken(token)
	if userID <= 0 {
		return item.Core{}, errors.New("data not found")
	}

	res, err := is.qry.Update(userID, itemID, updateData)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "internal server error"
		}
		return item.Core{}, errors.New(msg)
	}

	return res, nil
}

func (is *itemSrv) Delete(token interface{}, itemID uint) error {
	userID, _ := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("data not found")
	}

	err := is.qry.Delete(userID, itemID)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "internal server error"
		}
		return errors.New(msg)
	}

	return nil
}