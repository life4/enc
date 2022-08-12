package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KeyPublic struct {
	cfg Config
}

func (cmd KeyPublic) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "public",
		Aliases: []string{"public", "p"},
		Args:    cobra.NoArgs,
		Short:   "Convert private key to public key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	return c
}

func (cmd KeyPublic) run() error {
	key, err := ReadKeyStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	key, err = key.ToPublic()
	if err != nil {
		return fmt.Errorf("cannot convert key: %v", err)
	}
	b, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("cannot serialize key: %v", err)
	}
	_, err = cmd.cfg.Write(b)
	return err
}
