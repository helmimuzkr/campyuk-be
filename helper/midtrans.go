package helper

import (
	"campyuk-api/config"
	"errors"
	"log"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type PaymentGateway interface {
	ChargeTransaction(orderID string, grossAmt int, bank string) (string, error)
}

type midtransCore struct {
	core coreapi.Client
}

func NewCoreMidtrans(cfg *config.AppConfig) PaymentGateway {
	c := coreapi.Client{}
	c.New(cfg.SERVER_KEY, midtrans.Sandbox)

	return midtransCore{core: c}
}

func (c midtransCore) ChargeTransaction(orderID string, grossAmt int, bank string) (string, error) {

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(grossAmt),
		},
		BankTransfer: &coreapi.BankTransferDetails{
			Bank: midtrans.Bank(bank),
		},
		CustomExpiry: &coreapi.CustomExpiry{
			ExpiryDuration: 1,
			Unit:           "day",
		},
	}

	response, errMidtrans := c.core.ChargeTransaction(chargeReq)
	if errMidtrans != nil {
		log.Println(errMidtrans)
		return "", errors.New("charge transaction failed due to internal server error")
	}

	return response.VaNumbers[0].VANumber, nil
}
