package cmd

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Decrypt struct {
	cfg      Config
	password string
}

func (d Decrypt) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "decrypt",
		Aliases: []string{"decode", "d"},
		Short:   "Decrypt the message",
		RunE: func(cmd *cobra.Command, args []string) error {
			return d.run()
		},
	}
	c.Flags().StringVar(&d.password, "password", "", "password to use")
	return c
}

func (cmd Decrypt) run() error {
	data, err := io.ReadAll(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	var message *crypto.PGPMessage
	if bytes.HasPrefix(data, []byte("-----BEGIN PGP MESSAGE-----")) {
		message, err = crypto.NewPGPMessageFromArmored(string(data))
		if err != nil {
			return fmt.Errorf("cannot unarmor the message: %v", err)
		}
	} else {
		message = crypto.NewPGPMessage(data)
	}
	decrypted, err := crypto.DecryptMessageWithPassword(message, []byte(cmd.password))
	if err != nil {
		return fmt.Errorf("cannot decrypt the message: %v", err)
	}
	_, err = cmd.cfg.Write(decrypted.GetBinary())
	return err
}
