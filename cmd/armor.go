package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Armor struct {
	cfg     Config
	comment string
}

func (cmd Armor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "armor",
		Aliases: []string{"a"},
		Args:    cobra.NoArgs,
		Short:   "Convert the message (or key) from binary to text",
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

func (cmd Armor) run() error {
	data, err := io.ReadAll(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message := crypto.NewPGPMessage(data)
	armored, err := message.GetArmoredWithCustomHeaders(cmd.comment, ArmorHeaderVersion)
	if err != nil {
		return fmt.Errorf("cannot armor the message: %v", err)
	}
	_, err = cmd.cfg.Write([]byte(armored))
	return err
}
