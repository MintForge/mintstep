package basecoin

// Coin is the currency on the chain
type Coin struct {
	Amount    uint64 `json:"amount"`
	Denom     string `json:"denom"`
	percision uint64
}

func NewCoin(amount uint64) Coin {
	return Coin{Amount: amount, Denom: "BaseCoin", percision: 1000}
}
