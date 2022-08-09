package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

// Input (stdin):
// + binary file
// + text file
// + plain text
//
// Keys:
// + password
// + key (name, path, binary, or text)
// + key with passphrase
//
// Output (stdout):
// + binary
//
// Armoring should be a separate command.

type Encrypt struct {
	Stdout   io.Writer
	Stdin    io.Reader
	password string
}

func (e Encrypt) Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the message",
		RunE: func(cmd *cobra.Command, args []string) error {
			return e.run()
		},
	}
	c.Flags().StringVar(&e.password, "password", "", "password to use")
	return c
}

func (e Encrypt) run() error {
	data, err := io.ReadAll(e.Stdin)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	message := crypto.NewPlainMessage(data)
	encrypted, err := crypto.EncryptMessageWithPassword(message, []byte(e.password))
	if err != nil {
		return fmt.Errorf("cannot encrypt the message: %v", err)
	}
	_, err = e.Stdout.Write(encrypted.GetBinary())
	return err
}
