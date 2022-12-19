package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type SigVerify struct {
	cfg       Config
	password  string
	key       string
	signature string
}

func (cmd SigVerify) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "verify",
		Aliases: []string{"validate", "check", "v"},
		Args:    cobra.NoArgs,
		Short:   "Validate the message using signature",
		Example: "cat encrypted.bin | enc sig verify --key public.key --signature message.sig",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.signature, "signature", "s", "", "path to the signature")
	c.Flags().StringVarP(&cmd.password, "password", "p", "", "password to use to unlock the key")
	c.Flags().StringVarP(&cmd.key, "key", "k", "", "path to the key to use")
	Must(c.MarkFlagRequired("signature"))
	Must(c.MarkFlagFilename("signature"))
	Must(c.MarkFlagRequired("key"))
	Must(c.MarkFlagFilename("key"))
	return c
}

func (cmd SigVerify) run() error {
	message, err := ReadPlainMessageStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read message: %v", err)
	}

	f, err := os.Open(cmd.signature)
	if err != nil {
		return fmt.Errorf("open signature file: %v", err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read signature file: %v", err)
	}
	signature := crypto.NewPGPSignature(data)

	key, err := ReadKeyFile(cmd.key)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	if cmd.password != "" {
		key, err = key.Unlock([]byte(cmd.password))
		if err != nil {
			return fmt.Errorf("cannot unlock key: %v", err)
		}
	}
	keyring, err := crypto.NewKeyRing(key)
	if err != nil {
		return fmt.Errorf("cannot create keyring: %v", err)
	}
	err = keyring.VerifyDetached(message, signature, crypto.GetUnixTime())
	if err != nil {
		return err
	}
	fmt.Println("valid ✅")
	return nil
}
