package mintcoin

import "fmt"

// Coin is the currency on the chain
type Coin struct {
	Amount    uint64 `json:"amount"`
	Denom     string `json:"denom"`
	percision uint64
}

func NewCoin(amount uint64) Coin {
	return Coin{Amount: amount, Denom: "BaseCoin", percision: 1000}
}

func (c *Coin) Add(coin Coin) {
	c.Amount += coin.Amount
}

func (c *Coin) Sub(coin Coin) {
	c.Amount -= coin.Amount
}

func (c Coin) Bytes() []byte {
	return []byte(fmt.Sprintf("%v", c))
}
