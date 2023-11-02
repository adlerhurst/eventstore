package client

import (
	"encoding/json"
	"os"

	"github.com/adlerhurst/eventstore/service/internal/api/eventstore/v1alpha"
	"github.com/spf13/cobra"
)

var (
	payloadPath string
	pushCmd     = &cobra.Command{
		Use:   "push",
		Short: "calls the push method",
		Run:   push,
	}
)

func init() {
	pushCmd.PersistentFlags().StringVar(&payloadPath, "path", "", "path to the payload of the call")
	_ = pushCmd.MarkFlagFilename("path", ".json")
}

func push(cmd *cobra.Command, args []string) {
	req, err := readPayload(args)
	if err != nil {
		config.Logger.Error("failed to read payload", "cause", err)
		os.Exit(1)
	}

	_, err = client.Push(cmd.Context(), req)
	if err != nil {
		config.Logger.Error("failed to push", "cause", err)
		os.Exit(1)
	}
}

func readPayload(args []string) (req *eventstorev1alpha.PushRequest, err error) {
	if len(payloadPath) > 0 {
		return readPayloadFromFlag()
	}
	return nil, nil
}

func readPayloadFromFlag() (req *eventstorev1alpha.PushRequest, err error) {
	payload, err := os.ReadFile(payloadPath)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}

	return req, nil
}
