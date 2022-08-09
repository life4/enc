package cmd

import (
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/spf13/cobra"
)

type cfgEncrypt struct {
	Path string
}

func cmdEncrypt(w io.Writer) *cobra.Command {
	cfg := cfgEncrypt{}
	c := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the message",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := io.Open(c.Path)
			data, err := io.ReadAll(r)
			if err != nil {
				return fmt.Errorf("cannot read file: %v", err)
			}
			message := crypto.NewPlainMessageFromFile(data)
			encrypted, err := helper.EncryptMessageWithPassword(message, password)
		},
	}
	c.Flags().StringVarP()
	return c
}
