package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KeyLock struct {
	cfg      Config
	password string
}

func (cmd KeyLock) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "lock",
		Aliases: []string{"l"},
		Args:    cobra.NoArgs,
		Short:   "Protect the key with password",
		Example: "cat private.key | enc key lock --password 'my pass' > locked.key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.password, "password", "p", "", "password to use")
	Must(c.MarkFlagRequired("password"))
	return c
}

func (cmd KeyLock) run() error {
	key, err := ReadKeyStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("read key: %v", err)
	}
	key, err = key.Lock([]byte(cmd.password))
	if err != nil {
		return fmt.Errorf("lock key: %v", err)
	}
	b, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("serialize key: %v", err)
	}
	_, err = cmd.cfg.Write(b)
	return err
}
