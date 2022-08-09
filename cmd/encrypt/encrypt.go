package cmdencrypt

import (
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/spf13/cobra"
)

// Input (stdin?):
// + binary file
// + text file
// + plain text
//
// Keys:
// + password
// + key (name, path, binary, or text)
// + key with passphrase
//
// Output:
// + binary
// + text (armored)
//
// Armoring should be a separate command.

type Config struct {
	Path string
}

func Command(w io.Writer) *cobra.Command {
	// cfg := Config{}
	c := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the message",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("cannot read from stdin: %v", err)
			}
			message := crypto.NewPlainMessage(data)
			encrypted, err := helper.EncryptMessageWithPassword(message, password)
		},
	}
	return c
}
