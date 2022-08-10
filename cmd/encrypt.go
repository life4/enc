package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

// Keys:
// + password
// + key (name, path, binary, or text)
// + key with passphrase

type Encrypt struct {
	cfg      Config
	password string
}

func (e Encrypt) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "encrypt",
		Aliases: []string{"encode", "e"},
		Short:   "Encrypt the message",
		RunE: func(cmd *cobra.Command, args []string) error {
			return e.run()
		},
	}
	c.Flags().StringVar(&e.password, "password", "", "password to use")
	return c
}

func (cmd Encrypt) run() error {
	data, err := io.ReadAll(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message := crypto.NewPlainMessage(data)
	encrypted, err := crypto.EncryptMessageWithPassword(message, []byte(cmd.password))
	if err != nil {
		return fmt.Errorf("cannot encrypt the message: %v", err)
	}
	_, err = cmd.cfg.Write(encrypted.GetBinary())
	return err
}
