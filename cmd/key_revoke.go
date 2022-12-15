package cmd

import (
	"fmt"

	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/spf13/cobra"
)

type KeyRevoke struct {
	cfg    Config
	reason string
}

func (cmd KeyRevoke) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "revoke",
		Aliases: []string{"destroy", "r"},
		Args:    cobra.NoArgs,
		Short:   "Generate key revokation file",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.reason, "reason", "r", "", "a short explanation why the key is revoked")
	return c
}

func (cmd KeyRevoke) run() error {
	key, err := ReadKeyStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	entity := key.GetEntity()
	entity.RevokeKey(packet.NoReason, cmd.reason, nil)
	b, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("cannot revoke key: %v", err)
	}
	cmd.cfg.Write(b)
	return nil
}
