package cmd

import (
	"encoding/json"
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
	result := map[string]interface{}{
		"go_version": info.GoVersion,
		"revision":   getBuildKey(info, "vcs.revision"),
		"time":       getBuildKey(info, "vcs.time"),
		"os":         getBuildKey(info, "GOOS"),
		"arch":       getBuildKey(info, "GOARCH"),
		"now":        crypto.GetTime().Format(time.RFC3339),
	}

	b, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return fmt.Errorf("serialize JSON: %v", err)
	}
	_, err = cmd.cfg.Write(b)
	return err
}

func getBuildKey(info *debug.BuildInfo, key string) string {
	for _, setting := range info.Settings {
		if setting.Key == key {
			return setting.Value
		}
	}
	return ""
}
