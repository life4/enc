package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KeyUnlock struct {
	cfg      Config
	password string
}

func (cmd KeyUnlock) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "unlock",
		Aliases: []string{"u"},
		Args:    cobra.NoArgs,
		Short:   "Decode password-protected key",
		Example: "cat locked.key | enc key unlock --password 'my pass' > unlocked.key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.password, "password", "p", "", "password to use")
	c.MarkFlagRequired("password")
	return c
}

func (cmd KeyUnlock) run() error {
	key, err := ReadKeyStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("read key: %v", err)
	}
	key, err = key.Unlock([]byte(cmd.password))
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
