package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type SigID struct {
	cfg Config
}

func (cmd SigID) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "id",
		Aliases: []string{"ids", "info", "inspect", "i"},
		Args:    cobra.NoArgs,
		Short:   "Show ID (or IDs) of key used to create signature",
		Example: "cat message.sig | enc sig id",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	return c
}

func (cmd SigID) run() error {
	sig, err := ReadSigStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read signature: %v", err)
	}
	keyIDs, _ := sig.GetHexSignatureKeyIDs()
	for _, keyID := range keyIDs {
		fmt.Fprintln(cmd.cfg, keyID)
	}
	return nil
}
