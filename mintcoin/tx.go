package mintcoin

import (
	"bytes"
	"encoding/base64"

	"github.com/tendermint/tendermint/crypto"
)

type Tx struct {
	Type  byte   `json:"type"`
	RawTx string `json:"raw_tx"`
}

const (
	// TxTypeTransfer = byte(0x01)
	TxTypeCreate = byte(0x02)
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

func NewCreateTx(pub_key crypto.PubKey, signature []byte) *Tx {
	RegisterAmino(cdc)
	raw, _ := cdc.MarshalJSON(CreateTx{
		Address:   pub_key.Address(),
		PubKey:    pub_key,
		Signature: signature,
	})

	raw_str := base64.StdEncoding.EncodeToString(raw)

	return &Tx{
		Type:  TxTypeCreate,
		RawTx: raw_str,
	}
}
