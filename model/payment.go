package model

type Payment struct {
	Id              int
	OrderId         int    `json:"order_id"`
	Amount          int64  `json:"amount"`
	Status          string `json:"status"`
	Description     string `json:"description"`
	Currency        string `json:"currency"`
	CardName        string `json:"card_name"`
	CardType        string `json:"card_type"`
	CardNumber      string `json:"card_number"`
	CardExpiryMonth int    `json:"card_expiry_month"`
	CardExpiryYear  int    `json:"card_expiry_year"`
	Cvv             string `json:"cvv"`
}

func (p *Payment) Validate() bool {
	if p.Amount <= 0 || p.CardExpiryMonth <= 0 || p.CardExpiryMonth > 12 || len(p.Cvv) != 3 || len(p.CardNumber) != 16 {
		return false
	}

	return true
}

func (p *Payment) Masker() {
	p.CardNumber = maskLeft(p.CardNumber)
}

func maskLeft(s string) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-4; i++ {
		rs[i] = 'X'
	}
	return string(rs)
}
