package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/forlornfern/ffisow/internal"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "ffisow <iso> <device>",
		Short:   "write ISO-image to device",
		Example: "ffisow ~/Downloads/linux.iso /dev/sda",
		Args:    cobra.ExactArgs(2),
		PreRun: func(cmd *cobra.Command, args []string) {
			if verbose, _ := cmd.PersistentFlags().GetBool("verbose"); verbose {
				//
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
			if dirs, err := os.ReadDir("/sys/block"); err != nil {
				return err
			} else if !slices.ContainsFunc(dirs, func(dir os.DirEntry) bool {
				return dir.Name() == filepath.Base(args[1])
			}) {
				return fmt.Errorf("%q is not a block device", args[1])
			}
			if mounted, err := isMounted(args[1]); err != nil {
				return err
			} else if mounted {
				return fmt.Errorf("%q is mounted", args[1])
			}

			info, err := src.Stat()
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "All data on the %q will be overwritten\nAre you sure? [y/N] ", args[1])
			var input string
			fmt.Scanln(&input)
			fmt.Fprintf(os.Stderr, "\033[2A\033[0J")
			if strings.ToLower(input) != "y" {
				return nil
			}

			data, err := os.ReadFile(fmt.Sprintf("/sys/block/%s", filepath.Base(args[1])))
			if err != nil {
				return err
			}
			deviceSize, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
			if err != nil {
				return err
			}
			isoSize := info.Size()
			if isoSize > deviceSize*512 {
				return fmt.Errorf("%q size is to large for this block device", args[0])
			}

			pr := &internal.ProgressReader{
				Reader: src,
				Total:  isoSize,
			}

			bufSize, _ := cmd.Flags().GetInt64("buffer")

			buf := make([]byte, max(4, bufSize)*1024)
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
	rootCmd.Flags().Int64P("buffer", "b", 1024, "buffer size in KiB")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func isMounted(device string) (bool, error) {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return false, err
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && strings.HasPrefix(fields[0], device) {
			return true, nil
		}
	}
	return false, nil
}
