package model

// TBD

type Credit struct {
	CustomerId string     `json:"customerId"`
	OrderTotal float64    `json:"orderTotal"`
	CreditCard CreditCard `json:"creditCard"`
	OrderId    string     `json:"orderId"`
	ProductId  string     `json:"productId"`
}

type CreditCard struct {
	CardNumber string `json:"cardNumber"`
	Cvc        string `json:"cvc"`
	ExpiryDate string `json:"expiryDate"`
}
