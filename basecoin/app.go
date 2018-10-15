package basecoin

import (
	"bytes"
	"encoding/binary"
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
	result := ExecTx(app.state, txBytes)
	tags := []cmn.KVPair{}
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

func ExecTx(state *State, tx []byte) (res Result) {
	var key, value []byte
	parts := bytes.Split(tx, []byte("="))
	if len(parts) == 2 {
		key, value = parts[0], parts[1]
	} else {
		key, value = tx, tx
	}

	if bytes.Compare(key, []byte("address")) == 0 {
		account := NewAccount(value, 10000)
		SetAccount(state, value, *account)
	}
	res.Code = TypeOK
	return
}
