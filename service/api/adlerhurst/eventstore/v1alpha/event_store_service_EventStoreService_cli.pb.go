// Code generated by protoc-gen-cli-client. DO NOT EDIT.

package v1alpha

import (
	cli_client "github.com/adlerhurst/cli-client"
	cobra "github.com/spf13/cobra"
	grpc "google.golang.org/grpc"
	os "os"
)

var (
	ClientCmd = &cobra.Command{
		Use:                "client",
		Short:              ``,
		Long:               ``,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		DisableFlagParsing: true,
	}
)
var (
	_ClientPushCmdRequest = &PushRequestFlag{PushRequest: new(PushRequest)}
	ClientPushCmd         = &cobra.Command{
		Use:                "push",
		Short:              ``,
		Long:               ``,
		Run:                runClientPushCmd,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		DisableFlagParsing: true,
	}
	ClientPushCmdCallOptions []grpc.CallOption
)

func init() {
	ClientCmd.AddCommand(ClientPushCmd)
	ClientPushCmd.PreRun = func(cmd *cobra.Command, args []string) {
		ClientPushCmd.Flags().Parse(args)
		_ClientPushCmdRequest.AddFlags(ClientPushCmd.Flags())
		if ClientPushCmd.Flag("help").Changed {
			ClientPushCmd.Help()
			os.Exit(0)
		}
		_ClientPushCmdRequest.ParseFlags(cmd.Flags(), args)
	}
}

func runClientPushCmd(cmd *cobra.Command, args []string) {
	conn := cli_client.Connection(cmd.Context())
	client := NewEventStoreServiceClient(conn)

	res, err := client.Push(cmd.Context(), _ClientPushCmdRequest.PushRequest, ClientPushCmdCallOptions...)
	if err != nil {
		cli_client.Logger().Error("unable to Push", "cause", err)
		os.Exit(1)
	}
	cli_client.Logger().Info("🎉 request succeeded", "result", res)
}

var (
	_ClientFilterCmdRequest = &FilterRequestFlag{FilterRequest: new(FilterRequest)}
	ClientFilterCmd         = &cobra.Command{
		Use:                "filter",
		Short:              ``,
		Long:               ``,
		Run:                runClientFilterCmd,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		DisableFlagParsing: true,
	}
	ClientFilterCmdCallOptions []grpc.CallOption
)

func init() {
	ClientCmd.AddCommand(ClientFilterCmd)
	ClientFilterCmd.PreRun = func(cmd *cobra.Command, args []string) {
		ClientFilterCmd.Flags().Parse(args)
		_ClientFilterCmdRequest.AddFlags(ClientFilterCmd.Flags())
		if ClientFilterCmd.Flag("help").Changed {
			ClientFilterCmd.Help()
			os.Exit(0)
		}
		_ClientFilterCmdRequest.ParseFlags(cmd.Flags(), args)
	}
}

func runClientFilterCmd(cmd *cobra.Command, args []string) {
	conn := cli_client.Connection(cmd.Context())
	client := NewEventStoreServiceClient(conn)

	res, err := client.Filter(cmd.Context(), _ClientFilterCmdRequest.FilterRequest, ClientFilterCmdCallOptions...)
	if err != nil {
		cli_client.Logger().Error("unable to Filter", "cause", err)
		os.Exit(1)
	}
	defer res.CloseSend()
	resp, err := res.Recv()
	if err != nil {
		cli_client.Logger().Error("receive failed", "cause", err)
		os.Exit(1)
	}
	cli_client.Logger().Info("🎉 request succeeded", "result", resp)

}
