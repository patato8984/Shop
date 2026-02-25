package dto

import (
	"context"
	"fmt"
	"time"
)

type BankApi struct{}

func NewBankApi() *BankApi {
	return &BankApi{}
}
func (b BankApi) CreateTransaction(ctx context.Context, idOrder int, amount float64) (string, string, error) {
	time.Sleep(4 * time.Second)
	payment_id := fmt.Sprintf("pay_fake:%d_%d", idOrder, time.Now().Unix())
	url := "https://mock-bank.com/pay/" + payment_id
	return url, payment_id, nil
}
