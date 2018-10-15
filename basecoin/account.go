package basecoin

import (
	"encoding/json"
	"fmt"
)

type Account struct {
	Address []byte `json:"address"`
	Coin    Coin   `json:"coin"`
}

func (acc *Account) String() string {
	if acc == nil {
		return "invalid account"
	}
	return fmt.Sprintf("account %+v\n", *acc)
}

func (acc *Account) Add(coin Coin) {
	acc.Coin.Amount += coin.Amount
}

func (acc *Account) Sub(coin Coin) {
	acc.Coin.Amount -= coin.Amount
}

func GetAccount(state *State, addr []byte) (*Account, Result) {
	var res Result
	db := state.GetDB()
	data := db.Get(addr)
	if len(data) == 0 {
		res.Code = TypeAddressInvalid
		return nil, res
	}
	var acc Account
	err := json.Unmarshal(data, &acc)
	if err != nil {
		res.Code = TypeJsonEncodingError
		return nil, res
	}
	res.Code = TypeOK
	return &acc, res
}

func SetAccount(state *State, addr []byte, acc Account) {
	db := state.GetDB()
	accBytes, err := json.Marshal(acc)
	if err != nil {
		panic(err)
	}
	db.Set(addr, accBytes)
}

func NewAccount(address []byte, amount uint64) *Account {
	return &Account{
		Address: address,
		Coin:    NewCoin(amount),
	}
}
