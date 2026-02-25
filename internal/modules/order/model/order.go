package order_model

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
)

type Order struct {
	Id                 int           `json:"id,omitempty"`
	User_id            int           `json:"user_id,omitempty"`
	Address            string        `json:"address,omitempty"`
	Messages           string        `json:"messages,omitempty"`
	Status             string        `json:"status,omitempty"`
	Total_amount       int           `json:"total_amount,omitempty"`
	Update_at          sql.NullTime  `json:"update_at,omitempty"`
	Created_at         sql.NullTime  `json:"created_at,omitempty"`
	Price_all_products float64       `json:"Price_all_products,omitempty"`
	Items              []Order_items `json:"Items,omitempty"`
}
type Order_items struct {
	Id                int     `json:"id_orderItems"`
	Id_skus           int     `json:"id_skus"`
	Quantity          int     `json:"quantity"`
	Price_at_purchase float64 `json:"price_at_purchase"`
}
type BankApiRequest struct {
	PaymentId string `json:"payment_id"`
	Signature string `json:"hash"`
}
type BankTransactionResponse struct {
	Status  string `json:"status"`
	BankURL string `json:"bankURL"`
}

func (b BankApiRequest) HashComparison(payment_id, sharedSecret string) error {
	data := payment_id + sharedSecret
	hash := sha256.Sum256([]byte(data))
	myHashHex := hex.EncodeToString(hash[:])
	if b.Signature != myHashHex {
		return ErrHash
	}
	return nil
}

var (
	ErrHeader            = errors.New("header error")
	ErrQuantityItemsCart = errors.New("there are no items in the cart")
	ErrUrl               = errors.New("incorrect URL")
	ErrJson              = errors.New("error json")
	ErrHash              = errors.New("error hesh")
	ErrStatusPey         = errors.New("the order has already been paid for")
	ErrStatusCancelled   = errors.New("the order is cancelled ")
	ErrOrderNotFound     = errors.New("order not found")
	ErrStock             = errors.New("the product is out of stock")
	ErrCreateTransaction = errors.New("error bank api")
)
