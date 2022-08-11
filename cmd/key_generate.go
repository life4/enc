package cmd

import (
	"errors"
	"fmt"
	"os/user"

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

func (cmd KeyGenerate) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"create", "g"},
		Args:    cobra.NoArgs,
		Short:   "Generate new private key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVar(&cmd.name, "name", "", "your full name")
	c.Flags().StringVar(&cmd.email, "email", "", "your email address")
	c.Flags().StringVar(&cmd.ktype, "type", "rsa", "type of the key")
	c.Flags().IntVar(&cmd.bits, "bits", 4096, "size of the key in bits")
	c.Flags().BoolVar(&cmd.armor, "armor", false, "armor the key")
	return c
}

func (cmd KeyGenerate) run() error {
	username := cmd.Username()
	if username == "" {
		return errors.New("--name is required")
	}
	key, err := crypto.GenerateKey(username, cmd.email, cmd.ktype, cmd.bits)
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

func (cmd KeyGenerate) Username() string {
	if cmd.name != "" {
		return cmd.name
	}
	user, err := user.Current()
	if err == nil {
		if user.Name != "" {
			return user.Name
		}
		return user.Username
	}
	return ""
}
