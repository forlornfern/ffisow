package cmd

import (
	"io"
	"os"

	"charm.land/log/v2"
	"github.com/forlornfern/ffisow/internal"
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
			src, err := os.Open(args[0])
			if err != nil {
				return err
			}
			dst, err := os.OpenFile(args[1], os.O_WRONLY|os.O_SYNC, 0)
			if err != nil {
				return err
			}
			defer src.Close()
			defer dst.Close()
			info, err := src.Stat()
			if err != nil {
				return err
			}
			size := info.Size()

			pr := &internal.ProgressReader{
				Reader: src,
				Total:  size,
			}

			buf := make([]byte, 5*1024*1024)
			_, err = io.CopyBuffer(dst, pr, buf)
			if err != nil {
				return err
			}
			return dst.Sync()
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
