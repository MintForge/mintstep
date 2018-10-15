package basecoin

import (
	dbm "github.com/tendermint/tendermint/libs/db"
)

type State struct {
	db      dbm.DB
	Height  int64 `json:"height"`
	TxCount int64 `json:"tx_count"`
}

func (state State) GetDB() dbm.DB {
	return state.db
}

func InitState(db dbm.DB) *State {
	return &State{db: db, Height: 0, TxCount: 0}
}
