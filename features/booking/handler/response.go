package handler

type bookingResponse struct {
	ID            uint    `json:"id"`
	Ticket        string  `json:"ticket"`
	Title         string  `json:"camp_title"`
	Image         string  `json:"camp_image"`
	Latitude      float64 `json:"camp_latitude"`
	Longitude     float64 `json:"camp_longitude"`
	Address       string  `json:"camp_address"`
	City          string  `json:"camp_city"`
	CampPrice     string  `json:"camp_price,omitempty"`
	CheckIn       string  `json:"check_in"`
	CheckOut      string  `json:"check_out"`
	Guest         int     `json:"guest,omitempty"`
	CampCost      int     `json:"camp_cost,omitempty"`
	TotalPrice    int     `json:"total_price"`
	Status        string  `json:"status"`
	BookingDate   string  `json:"booking_date"`
	Bank          string  `json:"bank"`
	VirtualNumber string  `json:"virtual_number"`
}

type itemResponse struct {
	ID       uint   `json:"item_id,omitempty"`
	Name     string `json:"name"`
	Price    int    `json:"rent_price"`
	Quantity int    `json:"quantity"`
	RentCost int    `json:"rent_cost"`
}

type bookingDetailResponse struct {
	bookingResponse
	Items []itemResponse `json:"items"`
}

