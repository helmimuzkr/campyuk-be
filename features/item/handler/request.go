package handler

import "campyuk-api/features/item"

type AddItemRequest struct {
	CampID uint   `json:"camp_id" form:"camp_id"`
	Name   string `json:"name" form:"name"`
	Stock  int    `json:"stock" form:"stock"`
	Price  int    `json:"price" form:"price"`
}

type UpdateItemRequest struct {
	Name  string `json:"name" form:"name"`
	Stock int    `json:"stock" form:"stock"`
	Price int    `json:"price" form:"price"`
}

func RequestToCore(dataCart interface{}) *item.Core {
	res := item.Core{}
	switch dataCart.(type) {
	case AddItemRequest:
		cnv := dataCart.(AddItemRequest)
		res.CampID = int(cnv.CampID)
		res.Name = cnv.Name
		res.Stock = cnv.Stock
		res.Price = cnv.Price
	case UpdateItemRequest:
		cnv := dataCart.(UpdateItemRequest)
		res.Name = cnv.Name
		res.Stock = cnv.Stock
		res.Price = cnv.Price
	default:
		return nil
	}
	return &res

}
