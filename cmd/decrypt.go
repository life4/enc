package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type Decrypt struct {
	Stdout   io.Writer
	Stdin    io.Reader
	password string
}

func (d Decrypt) Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt the message",
		RunE: func(cmd *cobra.Command, args []string) error {
			return d.run()
		},
	}
	c.Flags().StringVar(&d.password, "password", "", "password to use")
	return c
}

func (e Decrypt) run() error {
	data, err := io.ReadAll(e.Stdin)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message := crypto.NewPGPMessage(data)
	decrypted, err := crypto.DecryptMessageWithPassword(message, []byte(e.password))
	_, err = e.Stdout.Write(decrypted.GetBinary())
	return err
}
