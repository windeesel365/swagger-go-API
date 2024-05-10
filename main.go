package main

type Shopper struct {
	Username   string `json:"username" gorm:"primaryKey"`
	FullName   string `json:"fullName"`
	Email      string `json:"email"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zipCode"`
	DateJoined string `json:"dateJoined"`
}

type ShoppersResponse struct {
	Shoppers []Shopper `json:"shoppers"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {

}
