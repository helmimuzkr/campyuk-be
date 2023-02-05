package handler

type bookingResponse struct {
	CampID     uint   `json:"camp_id"`
	Email      string `json:"email"` // Email guest
	Title      string `json:"title"`
	Image      string `json:"image"`
	Address    string `json:"address"`
	City       string `json:"city"`
	CampPrice  string `json:"camp_price"`
	CheckIn    string `json:"check_in"`
	CheckOut   string `json:"check_out"`
	Guest      int    `json:"guest"`
	CampCost   int    `json:"camp_cost"`
	TotalPrice int    `json:"total_price"`
}

type itemResponse struct {
	ID       uint `json:"item_id"`
	Name     string
	Price    int
	Quantity int `json:"quantity"`
	RentCost int `json:"rent_cost"`
}
