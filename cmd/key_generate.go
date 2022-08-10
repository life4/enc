package cmd

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/armor"
	"github.com/ProtonMail/gopenpgp/v2/constants"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type KeyGenerate struct {
	cfg   Config
	name  string
	email string
	ktype string
	bits  int
	armor bool
}

func (g KeyGenerate) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"create", "g"},
		Short:   "Generate new private key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.run()
		},
	}
	c.Flags().StringVar(&g.name, "name", "", "your full name")
	c.Flags().StringVar(&g.email, "email", "", "your email address")
	c.Flags().StringVar(&g.ktype, "type", "rsa", "type of the key")
	c.Flags().IntVar(&g.bits, "bits", 4096, "size of the key in bits")
	c.Flags().BoolVar(&g.armor, "armor", false, "armor the key")
	return c
}

func (cmd KeyGenerate) run() error {
	key, err := crypto.GenerateKey(cmd.name, cmd.email, cmd.ktype, cmd.bits)
	if err != nil {
		return fmt.Errorf("cannot generate key: %v", err)
	}
	b, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("cannot serialize key: %v", err)
	}
	if cmd.armor {
		s, err := armor.ArmorWithTypeAndCustomHeaders(
			b, constants.PrivateKeyHeader,
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
