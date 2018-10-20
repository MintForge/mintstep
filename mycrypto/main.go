package main

import (
	"fmt"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Account struct {
	PrivKey crypto.PrivKey `json:"priv_key"`
	Name    string         `json:"name"`
}

var cdc = amino.NewCodec()

func init() {
	RegisterAmino(cdc)
}

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

func main() {
	// privKey := secp256k1.GenPrivKey()
	// pubKey := privKey.PubKey()
	// var a crypto.PubKey
	// var b crypto.PrivKey
	// a = pubKey
	// b = privKey
	// msg := []byte("haha")
	// signature, _ := privKey.Sign(msg)
	// is_equal := pubKey.VerifyBytes(msg, signature)
	// fmt.Println(a)
	// fmt.Println(b)
	// fmt.Println(is_equal)
	account := Account{PrivKey: secp256k1.GenPrivKey(), Name: "haha"}
	// privKey := secp256k1.GenPrivKey()
	// jsonBytes, _ := cdc.MarshalJSON(privKey)
	jsonBytes, _ := cdc.MarshalJSON(account)
	// var p crypto.PrivKey
	var a Account
	cdc.UnmarshalJSON(jsonBytes, &a)
	fmt.Println(a.PrivKey)
}
