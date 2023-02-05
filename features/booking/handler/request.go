package handler

type bookingRequest struct {
	CampID     uint          `json:"camp_id" form:"camp_id"`
	CheckIn    string        `json:"check_in" form:"check_in"`
	CheckOut   string        `json:"check_out" form:"check_out"`
	Guest      int           `json:"guest" form:"guest"`
	CampCost   int           `json:"camp_cost" form:"camp_cost"`
	Items      []itemRequest `json:"items" form:"items"`
	TotalPrice int           `json:"total_price" form:"total_price"`
	Bank       string        `json:"bank" form:"bank"`
}

type itemRequest struct {
	ID       uint `json:"item_id" form:"item_id"`
	Quantity int  `json:"quantity" form:"quantity"`
	RentCost int  `json:"rent_cost" form:"rent_cost"`
}
