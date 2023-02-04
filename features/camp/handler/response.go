package handler

type campResponse struct {
	ID                 uint   `json:"id"`
	VerificationStatus string `json:"verification_status"`
	HostName           string `json:"host_name"`
	Title              string `json:"title"`
	Price              int    `json:"price"`
	Address            string `json:"address"`
	City               string `json:"city"`
	Distance           int    `json:"distance"`
	Document           string `json:"document,omitempty"`
	Image              string `json:"image,omitempty"`
}

type campDetailReponse struct {
	campResponse
	Images []campImage `json:"images"`
	Items  []campItem  `json:"items"`
}

type campItem struct {
	ID        uint   `json:"item_id"`
	Name      string `json:"name"`
	Stock     int    `json:"stock"`
	RentPrice int    `json:"rent_price"`
	ItemImage string `json:"item_image"`
}

type campImage struct {
	ID       uint   `json:"image_id"`
	ImageURL string `json:"image"`
}
