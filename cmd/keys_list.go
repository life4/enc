package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KeysList struct {
	cfg Config
}

func (cmd KeysList) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"all", "l"},
		Args:    cobra.NoArgs,
		Short:   "Show list of all keys in keyring",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	return c
}

func (cmd KeysList) run() error {
	keys, err := ReadKeys(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read keys: %v", err)
	}
	for _, key := range keys.GetKeys() {
		ident := key.GetEntity().PrimaryIdentity()
		var color string
		if key.IsExpired() || key.IsRevoked() {
			color = ColorRed
		} else if locked, _ := key.IsLocked(); locked {
			color = ColorGreen
		} else {
			color = ColorYellow
		}
		fmt.Fprintf(
			cmd.cfg, "%s%s%s %-35s %s\n",
			color, key.GetHexKeyID(), ColorEnd,
			ident.UserId.Email,
			ident.UserId.Name,
		)
	}
	return nil
}
