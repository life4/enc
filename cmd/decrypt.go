package cmd

import (
	"errors"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Decrypt struct {
	cfg      Config
	password string
	key      string
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
	c.Flags().StringVarP(&cmd.key, "key", "k", "", "password to use")
	c.MarkFlagFilename("key")
	return c
}

func (cmd Decrypt) run() error {
	message, err := ReadPGPMessageStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read message: %v", err)
	}

	var decrypted *crypto.PlainMessage
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
		if !key.IsPrivate() {
			return errors.New("public key cannot be used to decrypt")
		}
		keyring, err := crypto.NewKeyRing(key)
		if err != nil {
			return fmt.Errorf("cannot create keyring: %v", err)
		}
		decrypted, err = keyring.Decrypt(message, nil, 0)
		if err != nil {
			return fmt.Errorf("cannot decrypt message: %v", err)
		}
	} else if cmd.password != "" {
		decrypted, err = crypto.DecryptMessageWithPassword(message, []byte(cmd.password))
		if err != nil {
			return fmt.Errorf("cannot decrypt message: %v", err)
		}
	} else {
		return errors.New("a password or a key required")
	}
	_, err = cmd.cfg.Write(decrypted.GetBinary())
	return err
}
