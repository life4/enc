package cmd

import (
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Encrypt struct {
	cfg      Config
	password string
	key      string
}

func (cmd Encrypt) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "encrypt",
		Aliases: []string{"encode", "e", "rypt"},
		Args:    cobra.NoArgs,
		Short:   "Encrypt the message",
		Example: "echo 'my msg' | enc encrypt --password 'my pass' > encrypted.bin",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.password, "password", "p", "", "password to use")
	c.Flags().StringVarP(&cmd.key, "key", "k", "", "path to the key to use")
	c.MarkFlagFilename("key")
	return c
}

func (cmd Encrypt) run() error {
	message, err := ReadPlainMessageStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read message: %v", err)
	}
	var encrypted *crypto.PGPMessage
	if cmd.key != "" {
		key, err := ReadKeyFile(cmd.key)
		if err != nil {
			return fmt.Errorf("cannot read key: %v", err)
		}
		if cmd.password != "" {
			key, err = key.Unlock([]byte(cmd.password))
			if err != nil {
				return fmt.Errorf("cannot unlock key: %v", err)
			}
		}
		keyring, err := crypto.NewKeyRing(key)
		if err != nil {
			return fmt.Errorf("cannot create keyring: %v", err)
		}
		encrypted, err = keyring.Encrypt(message, nil)
	} else if cmd.password != "" {
		encrypted, err = crypto.EncryptMessageWithPassword(message, []byte(cmd.password))
	} else {
		return errors.New("a password or a key required")
	}
	if err != nil {
		return fmt.Errorf("cannot encrypt the message: %v", err)
	}
	_, err = cmd.cfg.Write(encrypted.GetBinary())
	return err
}
