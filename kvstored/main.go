package main

import (
	"github.com/MintForge/mintstep/kvstore"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	Address string = "tcp://0.0.0.0:26658"
	Abci    string = "socket"
)

func main() {
	var app types.Application
	app = kvstore.NewKVStoreApplication()

	srv, err := server.NewServer(Address, Abci, app)
	if err != nil {
		panic(err)
	}

	if err := srv.Start(); err != nil {
		panic(err)
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})
}
