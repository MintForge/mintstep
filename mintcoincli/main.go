package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MintForge/mintstep/mintcoin"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func loadConfig() *map[string]crypto.PrivKey {
	var filename = "accounts.json"
	var m map[string]crypto.PrivKey
	if _, err := os.Stat(filename); os.IsExist(err) {
		data, _ := ioutil.ReadFile(filename)
		_ = json.Unmarshal(data, &m)
	}
	return &m
}

var AccountMap = loadConfig()

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

// var TransferTxCmd = &cobra.Command{
// 	Use:   "transfer",
// 	Short: "transfer a tx",
// 	Long:  "transfer a tx",
// 	Args:  cobra.ExactArgs(3),
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		return cmdTransfer(cmd, args)
// 	},
// }

func cmdCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	var privKey crypto.PrivKey
	if _, ok := (*AccountMap)[name]; ok {
		privKey = (*AccountMap)[name]
	} else {
		privKey = secp256k1.GenPrivKey()
		(*AccountMap)[name] = privKey
	}
	signature, _ := privKey.Sign([]byte("create"))
	tx := mintcoin.NewCreateTx(privKey.PubKey(), signature)
	jsonBytes, _ := json.Marshal(tx)
	txBytes := base64.StdEncoding.EncodeToString(jsonBytes)
	url := fmt.Sprintf(`http://localhost:26657/broadcast_tx_commit?tx="%v"`, txBytes)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(string(body))
	return err
}

// func cmdTransfer(cmd *cobra.Command, args []string) error {
// 	sender := args[0]
// 	receiver := args[1]
// 	amount, _ := strconv.ParseInt(args[2], 10, 64)
// 	tx := basecoin.NewTransferTx([]byte(sender), 0, []byte(receiver), uint64(amount))
// 	jsonBytes, _ := json.Marshal(tx)
// 	txBytes := base64.StdEncoding.EncodeToString(jsonBytes)
// 	url := fmt.Sprintf(`http://localhost:26657/broadcast_tx_commit?tx="%v"`, txBytes)
// 	resp, err := http.Get(url)
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		// handle error
// 		panic(err)
// 	}

// 	fmt.Println(string(body))
// 	return err
// }

func addCommands() {
	RootCmd.AddCommand(CreateAccountCmd)
	// RootCmd.AddCommand(TransferTxCmd)
}

func Execute() error {
	addCommands()
	return RootCmd.Execute()
}

func main() {
	fmt.Println("+v", AccountMap)
	Execute()
	fmt.Println("+v", AccountMap)
	account_info, _ := json.Marshal(AccountMap)
	_ = ioutil.WriteFile("accounts.json", account_info, 0644)
}
