package pkg

import (
	"campyuk-api/config"
	"errors"
	"log"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type midtransAPI struct {
	core *coreapi.Client
}

func NewMidtrans(cfg *config.AppConfig) *midtransAPI {
	c := &coreapi.Client{}
	c.New(cfg.SERVER_KEY, midtrans.Sandbox)

	return &midtransAPI{core: c}
}

func (m *midtransAPI) ChargeTransaction(orderID string, grossAmt int, bank string) (string, error) {
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

	response, errMidtrans := m.core.ChargeTransaction(chargeReq)
	if errMidtrans != nil {
		log.Println(errMidtrans)
		return "", errors.New("charge transaction failed due to internal server error")
	}

	return response.VaNumbers[0].VANumber, nil
}
