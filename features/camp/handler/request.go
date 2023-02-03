package handler

type campRequest struct {
	Title       string  `json:"title" form:"title"`
	Price       int     `json:"price" form:"price"`
	Description string  `json:"description" form:"description"`
	Latitude    float64 `json:"latitude" form:"latitude"`
	Longitude   float64 `json:"longitude" form:"longitude"`
	Address     string  `json:"address" form:"address"`
	City        string  `json:"city" form:"city"`
	Distance    int     `json:"distance" form:"distance"`
}


