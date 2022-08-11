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
		Short:   "Revoke the key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVar(&cmd.reason, "reason", "", "a short explanation why the key is revoked")
	return c
}

func (cmd KeyRevoke) run() error {
	key, err := ReadKey(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	entity := key.GetEntity()
	entity.RevokeKey(packet.NoReason, "", nil)
	b, _ := key.Serialize()
	cmd.cfg.Write(b)
	return nil
}
