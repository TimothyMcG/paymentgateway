package db

import (
	"database/sql"
	"paymentgateway/model"
)

type db struct {
	db *sql.DB
}

func (d db) SaveUser(m model.Merchant) error {
	_, err := d.db.Exec("INSERT INTO processout.merchant(name) VALUES($1)", m.Name)
	if err != nil {
		return err
	}

	return nil
}

func (d db) VerifyMerchant(m model.Merchant) (int, bool, error) {
	err := d.db.QueryRow("SELECT id FROM processout.merchant WHERE name=$1",
		m.Name).Scan(&m.Id)
	if err != nil {
		return 0, false, err
	}

	return m.Id, true, nil
}

func (d db) SavePayment(p model.Payment) error {
	_, err := d.db.Exec("INSERT INTO processout.payments(id,card_number,card_type,card_name,cvv,card_expiry_month,card_expirt_year,amount,details, merchant_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", p.Id, p.CardNumber, p.CardType, p.CardName, p.Cvv, p.CardExpiryMonth, p.CardExpiryYear, p.Amount, p.Description, p.MerchantId)
	if err != nil {
		return err
	}

	return nil
}

func (d db) GetPayment(id int) (model.Payment, error) {
	var p model.Payment
	err := d.db.QueryRow("SELECT id, card_number,card_type,card_name,card_expiry_month,card_expirt_year,amount,details FROM processout.payments WHERE id=$1",
		id).Scan(&p.Id, &p.CardNumber, &p.CardType, &p.CardName, &p.CardExpiryMonth, &p.CardExpiryYear, &p.Amount, &p.Description)
	if err != nil {
		return model.Payment{}, err
	}

	return p, err
}
