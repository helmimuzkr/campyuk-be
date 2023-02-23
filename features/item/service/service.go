package service

import (
	"campyuk-api/features/item"
	"campyuk-api/pkg/helper"
	"errors"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

type itemSrv struct {
	qry item.ItemData
	vld *validator.Validate
}

func New(id item.ItemData, vld *validator.Validate) item.ItemService {
	return &itemSrv{
		qry: id,
		vld: vld,
	}
}

func (is *itemSrv) Add(token interface{}, campID uint, newItem item.Core) (item.Core, error) {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return item.Core{}, errors.New("access is denied due to invalid credential")
	}
	err := is.vld.Struct(&newItem)
	if err != nil {
		log.Println("err", err)
		msg := helper.ValidationErrorHandle(err)
		return item.Core{}, errors.New(msg)
	}

	res, err := is.qry.Add(userID, campID, newItem)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else if strings.Contains(err.Error(), "denied") {
			msg = err.Error()
		} else {
			msg = "internal server error"
		}
		log.Println(err)
		return item.Core{}, errors.New(msg)
	}

	return res, nil
}

func (is *itemSrv) Update(token interface{}, itemID uint, updateData item.Core) (item.Core, error) {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return item.Core{}, errors.New("access is denied due to invalid credential")
	}
	res, err := is.qry.Update(userID, itemID, updateData)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else if strings.Contains(err.Error(), "denied") {
			msg = err.Error()
		} else {
			msg = "internal server error"
		}
		log.Println(err)
		return item.Core{}, errors.New(msg)
	}

	return res, nil
}

func (is *itemSrv) Delete(token interface{}, itemID uint) error {
	userID, role := helper.ExtractToken(token)
	if role != "host" {
		return errors.New("access is denied due to invalid credential")
	}

	err := is.qry.Delete(userID, itemID)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else if strings.Contains(err.Error(), "denied") {
			msg = err.Error()
		} else {
			msg = "internal server error"
		}
		log.Println(err)
		return errors.New(msg)
	}

	return nil
}
