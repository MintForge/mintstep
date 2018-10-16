package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/MintForge/mintstep/basecoin"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: "the Client to basecoin",
}

var CreateAccountCmd = &cobra.Command{
	Use:   "create",
	Short: "create an account",
	Long:  "create an account from pubkey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdCreate(cmd, args)
	},
}

var TransferTxCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer a tx",
	Long:  "transfer a tx",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdTransfer(cmd, args)
	},
}

func cmdCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	tx := basecoin.NewCreateTx([]byte(name))
	txBytes, err := json.Marshal(tx)
	resp, err := http.Get(fmt.Sprintf(`http://localhost:26657/broadcast_tx_commit?tx="%v"`, txBytes))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(string(body))
	return err
}

func cmdTransfer(cmd *cobra.Command, args []string) error {
	sender := args[0]
	receiver := args[1]
	amount, _ := strconv.ParseInt(args[2], 10, 64)
	tx := basecoin.NewTransferTx([]byte(sender), 0, []byte(receiver), uint64(amount))
	txBytes, err := json.Marshal(tx)
	resp, err := http.Get(fmt.Sprintf(`http://localhost:26657/broadcast_tx_commit?tx="%v"`, txBytes))
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(string(body))
	return err
}

func addCommands() {
	RootCmd.AddCommand(CreateAccountCmd)
	RootCmd.AddCommand(TransferTxCmd)
}

func Execute() error {
	addCommands()
	return RootCmd.Execute()
}

func main() {
	Execute()
}
