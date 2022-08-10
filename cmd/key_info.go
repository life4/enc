package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KeyInfo struct {
	cfg Config
}

func (g KeyInfo) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "info",
		Aliases: []string{"inspect", "i"},
		Short:   "Show information about key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.run()
		},
	}
	return c
}

func (cmd KeyInfo) run() error {
	key, err := ReadKey(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	ident := key.GetEntity().PrimaryIdentity()
	w := cmd.cfg
	fmt.Fprintf(w, "{\n")

	// strings
	fmt.Fprintf(w, "  id: %#v,\n", key.GetHexKeyID())
	fmt.Fprintf(w, "  fingerprint: %#v,\n", key.GetFingerprint())
	fmt.Fprintf(w, "  name: %#v,\n", ident.UserId.Name)
	fmt.Fprintf(w, "  email: %#v,\n", ident.UserId.Email)
	fmt.Fprintf(w, "  comment: %#v,\n", ident.UserId.Comment)

	// flags
	fmt.Fprintf(w, "  is_private: %#v,\n", key.IsPrivate())
	fmt.Fprintf(w, "  is_expired: %#v,\n", key.IsExpired())
	fmt.Fprintf(w, "  is_revoked: %#v,\n", key.IsRevoked())
	fmt.Fprintf(w, "  can_verify: %#v,\n", key.CanVerify())
	fmt.Fprintf(w, "  can_encrypt: %#v,\n", key.CanEncrypt())

	fmt.Fprintf(w, "}\n")
	return nil
}
