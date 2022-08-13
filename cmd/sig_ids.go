package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type SigID struct {
	cfg Config
}

func (cmd SigID) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "id",
		Aliases: []string{"id", "info", "inspect", "i"},
		Args:    cobra.NoArgs,
		Short:   "Show ID (or IDs) of key used to create signature",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	return c
}

func (cmd SigID) run() error {
	if !cmd.cfg.HasStdin() {
		return errors.New("no signature passed into stdin")
	}
	data, err := io.ReadAll(cmd.cfg)
	if err != nil {
		return fmt.Errorf("read signature from stdin: %v", err)
	}
	signature := crypto.NewPGPSignature(data)
	keyIDs, _ := signature.GetHexSignatureKeyIDs()
	for _, keyID := range keyIDs {
		fmt.Fprintln(cmd.cfg, keyID)
	}
	return nil
}
