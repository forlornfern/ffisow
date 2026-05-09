package cmd

import (
	"os"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
)

var (
	logger  = log.NewWithOptions(os.Stderr, log.Options{Level: log.WarnLevel, ReportTimestamp: false})
	rootCmd = &cobra.Command{
		Use:           "ffisow <iso> <device>",
		Short:         "write ISO-image to device",
		Example:       "ffisow ~/Downloads/linux.iso /dev/sda1",
		Args:          cobra.ExactArgs(2),
		SilenceErrors: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			if verbose, _ := cmd.PersistentFlags().GetBool("verbose"); verbose {
				logger.SetLevel(log.InfoLevel)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose log")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
