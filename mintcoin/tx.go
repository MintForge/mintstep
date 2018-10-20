package mintcoin

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/tendermint/tendermint/crypto"
)

type Tx struct {
	Type  byte   `json:"type"`
	RawTx string `json:"raw_tx"`
}

const (
	TxTypeTransfer = byte(0x01)
	TxTypeCreate   = byte(0x02)
)

type TransferTx struct {
	Sender    crypto.Address `json:"sender"`
	PubKey    crypto.PubKey  `json:"pub_key"`
	Signature []byte         `json:"signature"`
	Coin      Coin           `json:"coin"`
	Sequence  uint64         `json:"sequence"`
	Receiver  crypto.Address `json:"receiver"`
}

type CreateTx struct {
	Address   crypto.Address `json:"address"`
	PubKey    crypto.PubKey  `json:"pub_key"`
	Signature []byte         `json:"signature"`
}

func (tx *CreateTx) Verify() bool {
	if bytes.Compare(tx.PubKey.Address(), tx.Address) != 0 {
		return false
	}
	return tx.PubKey.VerifyBytes([]byte("create"), tx.Signature)
}

func (tx *TransferTx) Verify() bool {
	if bytes.Compare(tx.PubKey.Address(), tx.Sender) != 0 {
		return false
	}
	container := []byte{}
	// msg := append(tx.Sender.Bytes(), tx.Receiver.Bytes()[:], tx.Coin.Bytes()[:]...)
	container = append(container, tx.Sender.Bytes()...)
	container = append(container, tx.Receiver.Bytes()...)
	container = append(container, tx.Coin.Bytes()...)
	return tx.PubKey.VerifyBytes(container, tx.Signature)
}

func NewTransferTx(pub_key crypto.PubKey, sequence uint64, signature []byte, receiver crypto.Address, amount uint64) *Tx {
	RegisterAmino(cdc)
	address := pub_key.Address()
	raw, err := cdc.MarshalJSON(TransferTx{
		Sender:    address,
		PubKey:    pub_key,
		Signature: signature,
		Coin:      NewCoin(amount),
		Sequence:  sequence,
		Receiver:  receiver,
	})
	fmt.Println(err)
	raw_str := base64.StdEncoding.EncodeToString(raw)

	return &Tx{
		Type:  TxTypeTransfer,
		RawTx: raw_str,
	}
}

func NewCreateTx(pub_key crypto.PubKey, signature []byte) *Tx {
	RegisterAmino(cdc)
	raw, err := cdc.MarshalJSON(CreateTx{
		Address:   pub_key.Address(),
		PubKey:    pub_key,
		Signature: signature,
	})

	fmt.Println(err)
	fmt.Println(raw)
	raw_str := base64.StdEncoding.EncodeToString(raw)

	fmt.Println(raw_str)

	return &Tx{
		Type:  TxTypeCreate,
		RawTx: raw_str,
	}
}
