package main

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func main() {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	var a crypto.PubKey
	var b crypto.PrivKey
	a = pubKey
	b = privKey
	msg := []byte("haha")
	signature, _ := privKey.Sign(msg)
	is_equal := pubKey.VerifyBytes(msg, signature)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(is_equal)
}
