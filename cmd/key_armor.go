package cmd

import (
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type KeyArmor struct {
	cfg     Config
	comment string
}

func (cmd KeyArmor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "armor",
		Aliases: []string{"a"},
		Args:    cobra.NoArgs,
		Short:   "Armor key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(
		&cmd.comment, "comment", "c", ArmorHeaderComment,
		"the comment to put into armored text",
	)
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
	s, err := key.ArmorWithCustomHeaders(cmd.comment, ArmorHeaderVersion)
	if err != nil {
		return fmt.Errorf("cannot armor key: %v", err)
	}
	_, err = cmd.cfg.Write([]byte(s))
	return err
}
