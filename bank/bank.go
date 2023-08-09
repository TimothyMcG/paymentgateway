package bank

var (
	count int
)

type IBankInterface interface {
	ProcessPayment(interface{}) bool
}

type BankRequest struct {
	CardName        string `json:"card_name"`
	CardType        string `json:"card_type"`
	CardNumber      string `json:"card_number"`
	CardExpiryMonth int    `json:"card_expiry_month"`
	CardExpiryYear  int    `json:"card_expiry_year"`
	Currency        string `json:"currency"`
	Amount          int64  `json:"amount"`
	Description     string `json:"description"`
}

func NewBankRequest() BankRequest {
	return BankRequest{}
}

func (b *BankRequest) ProcessPayment() (int, bool) {
	count++
	if count%2 == 0 {
		return 0, false
	}

	return count, true
}
