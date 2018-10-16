package basecoin

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var _ abci.Application = (*BaseCoinApplication)(nil)

type BaseCoinApplication struct {
	abci.BaseApplication
	state  *State
	logger *log.Logger
}

func NewBaseCoinApplication() *BaseCoinApplication {
	db := dbm.NewMemDB()
	state := InitState(db)
	logger := log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime)
	return &BaseCoinApplication{state: state, logger: logger}
}

func (app *BaseCoinApplication) Info(req abci.RequestInfo) abci.ResponseInfo {
	data := fmt.Sprintf("%+v", app.state)
	return abci.ResponseInfo{Data: data}
}

func (app *BaseCoinApplication) DeliverTx(txBytes []byte) abci.ResponseDeliverTx {

	tags := []cmn.KVPair{}
	result := ExecTx(app.state, txBytes)
	if result.IsErr() {
		return abci.ResponseDeliverTx{Code: TypeExecuteError, Tags: tags}
	}
	app.state.TxCount++
	return abci.ResponseDeliverTx{Code: TypeOK, Tags: tags}
}

func (app *BaseCoinApplication) CheckTx(txBytes []byte) abci.ResponseCheckTx {
	return abci.ResponseCheckTx{Code: TypeOK}
}

func (app *BaseCoinApplication) Commit() abci.ResponseCommit {
	app.state.Height++
	if app.state.TxCount == 0 {
		return abci.ResponseCommit{}
	}
	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, uint64(app.state.TxCount))
	return abci.ResponseCommit{Data: hash}
}

func (app *BaseCoinApplication) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		res.Log = "Query cannot be zero length"
		res.Code = TypeJsonEncodingError
		return
	}

	if req.Path == "account" {
		account, _ := GetAccount(app.state, req.Data)
		if account != nil {
			res.Value = account.Address
		} else {
			res.Log = "not found"
		}
	}

	return
}

func ExecTx(state *State, txBytes []byte) (res Result) {
	var raw_tx json.RawMessage
	tx := Tx{
		RawTx: &raw_tx,
	}
	if err := json.Unmarshal(txBytes, &tx); err != nil {
		res.Code = TypeJsonEncodingError
		return
	}
	switch tx.Type {
	case TxTypeCreate:
		var t CreateTx
		if err := json.Unmarshal(raw_tx, &t); err != nil {
			res.Code = TypeExecuteError
			return res
		}
		account := NewAccount(t.Name, 10000)
		SetAccount(state, t.Name, *account)
	case TxTypeTransfer:
		var t TransferTx
		if err := json.Unmarshal(raw_tx, &t); err != nil {
			res.Code = TypeExecuteError
			return res
		}
		sender, res := GetAccount(state, t.Sender)
		if res.IsErr() {
			res.Code = TypeAddressInvalid
			return res
		}
		receiver, res := GetAccount(state, t.Receiver)
		if res.IsErr() {
			res.Code = TypeAddressInvalid
			return res
		}
		sender.Sub(t.Coin)
		receiver.Add(t.Coin)
		res.Code = TypeOK
	}

	res.Code = TypeOK
	return res
}
