package cmd

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Version struct {
	cfg Config
}

func (cmd Version) Command() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print the version number",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
}

func (cmd Version) run() error {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("cannot read build info")
	}
	w := cmd.cfg
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
