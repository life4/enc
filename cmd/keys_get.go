package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

type KeysGet struct {
	cfg   Config
	query string
}

func (cmd KeysGet) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "get",
		Aliases: []string{"filter", "g", "export"},
		Short:   "Get a specific key from keyring",
		Example: "cat ~/.gnupg/pubring.gpg | enc keys get 514292cf25399377 > public.key",
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cmd.query = args[0]
			return cmd.run()
		},
	}
	return c
}

func (cmd KeysGet) run() error {
	keys, err := ReadKeys(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read keys: %v", err)
	}
	for _, key := range keys.GetKeys() {
		user := key.GetEntity().PrimaryIdentity().UserId
		if user.Email == cmd.query || key.GetHexKeyID() == cmd.query {
			b, err := key.Serialize()
			if err != nil {
				return fmt.Errorf("cannot serialize key: %v", err)
			}
			_, err = cmd.cfg.Write(b)
			return err
		}
	}
	return errors.New("key not found")
}
