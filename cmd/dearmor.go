package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

// TODO(gram): a better name
type Dearmor struct {
	Stdout io.Writer
	Stdin  io.Reader
}

func (d Dearmor) Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "dearmor",
		Short: "Convert the message (or key) from text to binary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return d.run()
		},
	}
	return c
}

func (d Dearmor) run() error {
	data, err := io.ReadAll(d.Stdin)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message, err := crypto.NewPGPMessageFromArmored(string(data))
	if err != nil {
		return fmt.Errorf("cannot de-armor the message: %v", err)
	}
	_, err = d.Stdout.Write(message.GetBinary())
	return err
}
