package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

// TODO(gram): a better name
type Dearmor struct {
	cfg Config
}

func (d Dearmor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "dearmor",
		Aliases: []string{"enarmor", "unarmor", "u"},
		Short:   "Convert the message (or key) from text to binary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return d.run()
		},
	}
	return c
}

func (d Dearmor) run() error {
	data, err := io.ReadAll(d.cfg)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message, err := crypto.NewPGPMessageFromArmored(string(data))
	if err != nil {
		return fmt.Errorf("cannot de-armor the message: %v", err)
	}
	_, err = d.cfg.Write(message.GetBinary())
	return err
}
