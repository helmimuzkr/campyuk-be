package handler

type listCampResponse struct {
	ID                 uint     `json:"id"`
	VerificationStatus string   `json:"verification_status"`
	HostName           string   `json:"host_name"`
	Title              string   `json:"title"`
	Price              int      `json:"price"`
	City               string   `json:"city"`
	Distance           int      `json:"distance"`
	Images             []string `json:"images"`
}

// 	"id": 1,
//   "verification_status": "accepted",
//   "host_name": "john",
//   "title": "Tanakita Camp",
//   "price": 100000,
//   "city": "Gotham City",
//   "distance": 100,
//   "image": "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fassets.pikiran-rakyat.com%2Fcrop%2F0x0%3A0x0%2Fx%2Fphoto%2F2020%2F04%2F04%2F2309824160.jpg&f=1&nofb=1&ipt=40c1a6e95ad3ea4a52e83708d27385d1436ef54e77c3afb84623998a0120f9eb&ipo=images"
