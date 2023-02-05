package handler

type campResponse struct {
	ID                 uint    `json:"id"`
	VerificationStatus string  `json:"verification_status"`
	HostName           string  `json:"host_name"`
	Title              string  `json:"title"`
	Price              int     `json:"price"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Address            string  `json:"address"`
	City               string  `json:"city"`
	Distance           int     `json:"distance"`
	Document           string  `json:"document,omitempty"`
	Image              string  `json:"image,omitempty"`
}

type campItem struct {
	ID        uint   `json:"item_id"`
	Name      string `json:"name"`
	Stock     int    `json:"stock"`
	RentPrice int    `json:"rent_price"`
}

type campImage struct {
	ID       uint   `json:"image_id"`
	ImageURL string `json:"image"`
}

type campDetailReponse struct {
	campResponse
	Images []campImage `json:"images"`
	Items  []campItem  `json:"items"`
}

type paginationResponse struct {
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`
	TotalRecord int `json:"total_rercord"`
	TotalPage   int `json:"total_page"`
}

type withPagination struct {
	Pagination paginationResponse `json:"pagination"`
	Data       []campResponse     `json:"data"`
	Message    string             `json:"message"`
}
