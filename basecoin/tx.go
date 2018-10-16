package basecoin

type Tx struct {
	Type  byte        `json:"type"`
	RawTx interface{} `json:"raw_tx"`
}

const (
	TxTypeTransfer = byte(0x01)
	TxTypeCreate   = byte(0x02)
)

type TransferTx struct {
	Sender   []byte `json:"sender"`
	Coin     Coin   `json:"coin"`
	Sequence uint64 `json:"sequence"`
	Receiver []byte `json:"receiver"`
}

type CreateTx struct {
	Name []byte `json:"name"`
}

func NewTransferTx(sender []byte, sequence uint64, receiver []byte, amount uint64) *Tx {
	return &Tx{
		Type: TxTypeTransfer,
		RawTx: TransferTx{
			Sender:   sender,
			Coin:     NewCoin(amount),
			Sequence: sequence,
			Receiver: receiver,
		},
	}
}

func NewCreateTx(name []byte) *Tx {
	return &Tx{
		Type: TxTypeCreate,
		RawTx: CreateTx{
			Name: name,
		},
	}
}
