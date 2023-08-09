package db

import "paymentgateway/model"

type Db interface {
	SaveUser(model.Merchant) error
	VerifyMerchant(model.Merchant) (int, bool, error)
	GetPayment(int) (model.Payment, error)
	SavePayment(model.Payment) error
}
