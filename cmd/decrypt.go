package cmd

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Decrypt struct {
	cfg      Config
	password string
}

func (cmd Decrypt) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "decrypt",
		Aliases: []string{"decode", "d"},
		Args:    cobra.NoArgs,
		Short:   "Decrypt the message",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.password, "password", "p", "", "password to use")
	return c
}

func (cmd Decrypt) run() error {
	message, err := ReadPGPMessageStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read message: %v", err)
	}
	decrypted, err := crypto.DecryptMessageWithPassword(message, []byte(cmd.password))
	if err != nil {
		return fmt.Errorf("cannot decrypt message: %v", err)
	}
	_, err = cmd.cfg.Write(decrypted.GetBinary())
	return err
}
