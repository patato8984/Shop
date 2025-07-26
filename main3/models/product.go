package models

type Product = struct {
	Name        string `json:"name"`
	Id          int    `json:"id"`
	Price       int    `json:"price"`
	Capacitance int    `json:"capacitance"`
	Hertz       int    `json:"hertz"`
	Status      string `json:"status"`
}
