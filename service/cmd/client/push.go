package client

import (
	"os"

	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	pushCmd = &cobra.Command{
		Use:                "push",
		Short:              "calls the push method",
		Run:                push,
		PreRun:             parsePushRequest,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		DisableFlagParsing: true,
	}
)

func init() {
	pushCmd.Flags().StringVar(&payloadPath, "path", "", "path to the payload of the call")
	err := pushCmd.MarkFlagFilename("path", ".json")
	if err != nil {
		config.Logger.Error("failed to mark flag filename", "cause", err)
		os.Exit(1)
	}
	pushCmd.Flags()
}

func push(cmd *cobra.Command, args []string) {
	if len(pushRequest.Aggregates) == 0 {
		config.Logger.Info("no valid aggregates provided, skip execution")
		return
	}
	_, err := client.Push(cmd.Context(), pushRequest)
	if err != nil {
		config.Logger.Error("failed to push", "cause", err)
		os.Exit(1)
	}
}

func readPayloadFromFlags(req *eventstorev1alpha.PushRequest) (err error) {
	payload, err := os.ReadFile(payloadPath)
	if err != nil {
		return err
	}

	if err = protojson.Unmarshal(payload, req); err != nil {
		return err
	}

	return nil
}
