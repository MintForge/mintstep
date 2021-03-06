package main

import (
	"github.com/MintForge/mintstep/basecoin"
	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	flagAddress string = "tcp://0.0.0.0:26658"
	flagAbci    string = "socket"
)

func main() {
	var app types.Application
	app = basecoin.NewBaseCoinApplication()

	srv, err := server.NewServer(flagAddress, flagAbci, app)
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
