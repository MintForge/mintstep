package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MintForge/mintstep/mintcoin"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var cdc = amino.NewCodec()

func RegisterAmino(cdc *amino.Codec) {
	// These are all written here instead of
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{},
		secp256k1.PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{},
		multisig.PubKeyMultisigThresholdAminoRoute, nil)

	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{},
		ed25519.PrivKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{},
		secp256k1.PrivKeyAminoRoute, nil)
}

func loadConfig() *map[string]crypto.PrivKey {
	RegisterAmino(cdc)
	var filename = "accounts.json"
	m := make(map[string]crypto.PrivKey)
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		data, _ := ioutil.ReadFile(filename)
		err = cdc.UnmarshalJSON(data, &m)
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
	jsonBytes, _ := cdc.MarshalJSON(tx)
	txBytes := base64.StdEncoding.EncodeToString(jsonBytes)
	url := fmt.Sprintf(`http://localhost:26657/broadcast_tx_commit?tx="%v"`, txBytes)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		panic(err)
	}

	fmt.Println(privKey.PubKey().Address())

	fmt.Println(string(body))
	return err
}

// func cmdTransfer(cmd *cobra.Command, args []string) error {
// }

func addCommands() {
	RootCmd.AddCommand(CreateAccountCmd)
}

func Execute() error {
	addCommands()
	return RootCmd.Execute()
}

func main() {
	Execute()
	account_info, _ := cdc.MarshalJSON(*AccountMap)
	_ = ioutil.WriteFile("accounts.json", account_info, 0644)
}
