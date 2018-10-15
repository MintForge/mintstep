package basecoin

type Tx struct {
	Sender   []byte `json:"sender"`
	Coin     Coin   `json:"coin"`
	Sequence uint64 `json:"sequence"`
	Receiver []byte `json:"receiver"`
}

func NewTx(sender []byte, sequence uint64, receiver []byte, coin Coin) *Tx {
	return &Tx{
		Sender:   sender,
		Coin:     coin,
		Sequence: sequence,
		Receiver: receiver,
	}
}
