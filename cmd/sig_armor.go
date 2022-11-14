package cmd

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/armor"
	"github.com/ProtonMail/gopenpgp/v2/constants"
	"github.com/spf13/cobra"
)

type SigArmor struct {
	cfg     Config
	comment string
}

func (cmd SigArmor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "armor",
		Aliases: []string{"a"},
		Args:    cobra.NoArgs,
		Short:   "Armor key",
		Example: "cat message.sig | enc sig armor > message-sig.txt",
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

func (cmd SigArmor) run() error {
	sig, err := ReadSigStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read signature: %v", err)
	}
	s, err := armor.ArmorWithType(sig.Data, constants.PGPSignatureHeader)
	if err != nil {
		return fmt.Errorf("cannot armor signature: %v", err)
	}
	_, err = cmd.cfg.Write([]byte(s))
	return err
}
