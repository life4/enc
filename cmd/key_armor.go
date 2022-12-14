package cmd

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/armor"
	"github.com/ProtonMail/gopenpgp/v2/constants"
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
		Example: "cat private.key | enc key armor > private-key.txt",
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
	key, err := ReadKeyStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	serialized, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("cannot serialize key: %v", err)
	}
	header := constants.PublicKeyHeader
	if key.IsPrivate() {
		header = constants.PrivateKeyHeader
	}
	armored, err := armor.ArmorWithTypeAndCustomHeaders(
		serialized, header, ArmorHeaderVersion, cmd.comment,
	)
	if err != nil {
		return fmt.Errorf("cannot armor key: %v", err)
	}
	_, err = cmd.cfg.Write([]byte(armored))
	return err
}
