package mintcoin

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"
)

type Account struct {
	Coin     Coin          `json:"coin"`
	PubKey   crypto.PubKey `json:"public_key"`
	Sequence int64         `json:"sequence"`
}

func NewAccount(pub_key crypto.PubKey, amount uint64) Account {
	return Account{
		Coin:     NewCoin(amount),
		PubKey:   pub_key,
		Sequence: 0,
	}
}

func (acc Account) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

func (acc Account) GetAddress() crypto.Address {
	return acc.PubKey.Address()
}

func (acc Account) GetSequence() int64 {
	return acc.Sequence
}

func (acc Account) String() string {
	return fmt.Sprintf("account:%v\n", acc.PubKey.Address())
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
	err := cdc.UnmarshalJSON(data, &acc)
	if err != nil {
		res.Code = TypeJsonEncodingError
		return nil, res
	}
	res.Code = TypeOK
	return &acc, res
}

func SetAccount(state *State, addr crypto.Address, acc Account) {
	db := state.GetDB()
	accBytes, err := cdc.MarshalJSON(acc)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db.Set([]byte(addr.String()), accBytes)
}

type PriviteAccount struct {
	Name    string         `json:"name"`
	Privkey crypto.PrivKey `json:"priv_key"`
	Account
}
