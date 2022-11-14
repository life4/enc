package cmd

import (
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type KeyDearmor struct {
	cfg Config
}

func (cmd KeyDearmor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "dearmor",
		Aliases: []string{"d", "enarmor", "unarmor"},
		Args:    cobra.NoArgs,
		Short:   "Dearmor key",
		Example: "cat private-key.txt | enc key dearmor > private.key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	return c
}

func (cmd KeyDearmor) run() error {
	if !cmd.cfg.HasStdin() {
		return errors.New("no key passed into stdin")
	}
	key, err := crypto.NewKeyFromArmoredReader(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	b, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("cannot serialize key: %v", err)
	}
	_, err = cmd.cfg.Write(b)
	return err
}
