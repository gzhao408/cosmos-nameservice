package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	app "github.com/gzhao408/cosmos-nameservice"
	nameservicecmd "github.com/cosmos/sdk-application-tutorial/x/nameservice/client/cli"
)

const storeAcc = "acc"

var (
	rootCmd = &cobra.Command{
		Use:   "nameservicecli",
		Short: "nameservice Client",
	}
	defaultCLIHome = os.ExpandEnv("$HOME/.nameservicecli")
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	rootCmd.AddCommand(client.ConfigCmd())
	rpc.AddCommands(rootCmd)

	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(queryCmd, cdc)
	queryCmd.AddCommand(client.LineBreak)
	queryCmd.AddCommand(client.GetCommands(
		authcmd.GetAccountCmd(storeAcc, cdc, authcmd.GetAccountDecoder(cdc)),
		nameservicecmd.GetCmdResolveName("nameservice", cdc),
		nameservicecmd.GetCmdWhois("nameservice", cdc),
	)...)

	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		nameservicecmd.GetCmdBuyName(cdc),
		nameservicecmd.GetCmdSetName(cdc),
	)...)

	rootCmd.AddCommand(
		queryCmd,
		txCmd,
		client.LineBreak,
	)

	rootCmd.AddCommand(
		keys.Commands(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "NS", defaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
