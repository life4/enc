package cmd

import (
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type KeyArmor struct {
	cfg Config
}

func (g KeyArmor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "armor",
		Aliases: []string{"a"},
		Args:    cobra.NoArgs,
		Short:   "Armor key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.run()
		},
	}
	return c
}

func (cmd KeyArmor) run() error {
	if !cmd.cfg.HasStdin() {
		return errors.New("no key passed into stdin")
	}
	key, err := crypto.NewKeyFromReader(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	s, err := key.ArmorWithCustomHeaders(ArmorHeaderComment, ArmorHeaderVersion)
	if err != nil {
		return fmt.Errorf("cannot armor key: %v", err)
	}
	_, err = cmd.cfg.Write([]byte(s))
	return err
}
