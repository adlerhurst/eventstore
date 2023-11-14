// Code generated by protoc-gen-go-cli. DO NOT EDIT.
//00:31:01

package eventstorev1alpha

import (
	cobra "github.com/spf13/cobra"
	pflag "github.com/spf13/pflag"
	os "os"
)

func UnmarshalClientPushRequest(cmd *cobra.Command, args []string) {
	set := pflag.NewFlagSet("request", pflag.ContinueOnError)

	new(PathField).AddFlag(set)

	cmd.Flags().AddFlagSet(set)
	cmd.DisableFlagParsing = false
	if err := cmd.ParseFlags(args); err != nil {
		DefaultConfig.Logger.Error("failed to parse flags", "cause", err)
		os.Exit(1)
	}
}
