package version

import (
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

func Command(w io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			return version(w)
		},
	}
}

func version(w io.Writer) error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("cannot read build info")
	}
	fmt.Fprintf(w, "{\n")
	fmt.Fprintf(w, "  go_version: %#v,\n", info.GoVersion)
	fmt.Fprintf(w, "  revision:   %#v,\n", getBuildKey(info, "vcs.revision"))
	fmt.Fprintf(w, "  time:       %#v,\n", getBuildKey(info, "vcs.time"))
	fmt.Fprintf(w, "  os:         %#v,\n", getBuildKey(info, "GOOS"))
	fmt.Fprintf(w, "  arch:       %#v,\n", getBuildKey(info, "GOARCH"))
	fmt.Fprintf(w, "  now:        %#v,\n", crypto.GetTime().Format(time.RFC3339))
	fmt.Fprintf(w, "}\n")
	return nil
}

func getBuildKey(info *debug.BuildInfo, key string) string {
	for _, setting := range info.Settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return ""
}
