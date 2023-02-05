package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func NewCoreMidtrans() coreapi.Client {
	c := coreapi.Client{}
	c.New(SERVER_KEY, midtrans.Sandbox)

	return c
}
