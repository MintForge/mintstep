package mintcoin

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var _ abci.Application = (*MintCoinApplication)(nil)

type MintCoinApplication struct {
	abci.BaseApplication
	state  *State
	logger *log.Logger
}

func NewMintCoinApplication() *MintCoinApplication {
	db := dbm.NewMemDB()
	state := InitState(db)
	logger := log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime)
	RegisterAmino(cdc)
	return &MintCoinApplication{state: state, logger: logger}
}

func (app *MintCoinApplication) Info(req abci.RequestInfo) abci.ResponseInfo {
	data := fmt.Sprintf("%+v", app.state)
	return abci.ResponseInfo{Data: data}
}

func (app *MintCoinApplication) DeliverTx(txBytes []byte) abci.ResponseDeliverTx {
	decoded, err := base64.StdEncoding.DecodeString(string(txBytes))
	tags := []cmn.KVPair{}
	if err != nil {
		return abci.ResponseDeliverTx{Code: TypeJsonEncodingError, Tags: tags}
	}
	result := ExecTx(app.state, decoded)
	if result.IsErr() {
		return abci.ResponseDeliverTx{Code: TypeExecuteError, Tags: tags}
	}
	app.state.TxCount++
	return abci.ResponseDeliverTx{Code: TypeOK, Tags: tags}
}

func (app *MintCoinApplication) CheckTx(txBytes []byte) abci.ResponseCheckTx {
	return abci.ResponseCheckTx{Code: TypeOK}
}

func (app *MintCoinApplication) Commit() abci.ResponseCommit {
	app.state.Height++
	if app.state.TxCount == 0 {
		return abci.ResponseCommit{}
	}
	hash := make([]byte, 8)
	binary.BigEndian.PutUint64(hash, uint64(app.state.TxCount))
	return abci.ResponseCommit{Data: hash}
}

func (app *MintCoinApplication) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		res.Log = "Query cannot be zero length"
		res.Code = TypeJsonEncodingError
		return
	}

	if req.Path == "account" {
		account, _ := GetAccount(app.state, req.Data)
		if account != nil {
			res.Log = account.String()
		} else {
			res.Log = "not found"
		}
	}

	return
}

func ExecTx(state *State, txBytes []byte) (res Result) {
	var tx Tx
	if err := json.Unmarshal(txBytes, &tx); err != nil {
		res.Code = TypeJsonEncodingError
		return
	}
	raw_str := tx.RawTx
	jsonBytes, _ := base64.StdEncoding.DecodeString(raw_str)
	switch tx.Type {
	case TxTypeCreate:
		var t CreateTx
		if err := cdc.UnmarshalJSON(jsonBytes, &t); err != nil {
			res.Code = TypeExecuteError
			return res
		}
		if t.Verify() {
			account := NewAccount(t.PubKey, 10000)
			SetAccount(state, t.PubKey.Address(), account)
		} else {
			res.Code = TypeVerifyError
			return res
		}

	}

	res.Code = TypeOK
	return res
}
