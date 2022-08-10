package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Armor struct {
	cfg Config
}

func (a Armor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "armor",
		Short: "Convert the message (or key) from binary to text",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.run()
		},
	}
	return c
}

func (e Armor) run() error {
	data, err := io.ReadAll(e.cfg)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message := crypto.NewPGPMessage(data)
	armored, err := message.GetArmored()
	if err != nil {
		return fmt.Errorf("cannot armor the message: %v", err)
	}
	_, err = e.cfg.Write([]byte(armored))
	return err
}
