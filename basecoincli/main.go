package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cli",
	Short: "the Client to basecoin",
}

var createAccountCmd = &cobra.Command{
	Use:   "create",
	Short: "create an account",
	Long:  "create an account from pubkey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdCreateAccount(cmd, args)
	},
}

func cmdCreateAccount(cmd *cobra.Command, args []string) error {
	address := args[0]
	resp, err := http.Get(fmt.Sprintf(`http://localhost:26657/broadcast_tx_commit?tx="address=%v"`, address))
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
	RootCmd.AddCommand(createAccountCmd)
}

func Execute() error {
	addCommands()
	return RootCmd.Execute()
}

func main() {
	Execute()
}
