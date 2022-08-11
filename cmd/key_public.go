package cmd

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/armor"
	"github.com/ProtonMail/gopenpgp/v2/constants"
	"github.com/spf13/cobra"
)

type KeyPublic struct {
	cfg   Config
	armor bool
}

func (g KeyPublic) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "public",
		Aliases: []string{"public", "p"},
		Args:    cobra.NoArgs,
		Short:   "Convert private key to public key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.run()
		},
	}
	c.Flags().BoolVar(&g.armor, "armor", false, "armor the key")
	return c
}

func (cmd KeyPublic) run() error {
	key, err := ReadKey(cmd.cfg)
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
	if cmd.armor {
		s, err := armor.ArmorWithTypeAndCustomHeaders(
			b, constants.PublicKeyHeader,
			ArmorHeaderVersion, ArmorHeaderComment,
		)
		if err != nil {
			return fmt.Errorf("cannot armor key: %v", err)
		}
		b = []byte(s)
	}
	_, err = cmd.cfg.Write(b)
	return err
}
