package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type SigCreate struct {
	cfg      Config
	password string
	key      string
}

func (cmd SigCreate) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "create",
		Aliases: []string{"sign", "generate", "c", "n"},
		Args:    cobra.NoArgs,
		Short:   "Sign the message",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.password, "password", "p", "", "password to use to unlock the key")
	c.Flags().StringVarP(&cmd.key, "key", "k", "", "path to the key to use")
	return c
}

func (cmd SigCreate) run() error {
	if !cmd.cfg.HasStdin() {
		return errors.New("no file passed into stdin")
	}
	data, err := io.ReadAll(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message := crypto.NewPlainMessage(data)
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
	signature, err := keyring.SignDetached(message)
	if err != nil {
		return fmt.Errorf("cannot encrypt the message: %v", err)
	}
	_, err = cmd.cfg.Write(signature.GetBinary())
	return err
}
