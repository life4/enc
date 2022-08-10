package cmd

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type KeyInfo struct {
	cfg Config
}

func (g KeyInfo) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "info",
		Aliases: []string{"inspect", "i"},
		Short:   "Show information about the key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.run()
		},
	}
	return c
}

func (cmd KeyInfo) run() error {
	data, err := io.ReadAll(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read from stdin: %v", err)
	}
	var key *crypto.Key
	isArmored := bytes.HasPrefix(data, []byte("-----BEGIN PGP PRIVATE KEY BLOCK-----"))
	if !isArmored {
		isArmored = bytes.HasPrefix(data, []byte("-----BEGIN PGP PUBLIC KEY BLOCK-----"))
	}
	if isArmored {
		key, err = crypto.NewKeyFromArmored(string(data))
		if err != nil {
			return fmt.Errorf("cannot unarmor the key: %v", err)
		}
	} else {
		key, err = crypto.NewKey(data)
		if err != nil {
			return fmt.Errorf("cannot parse the key: %v", err)
		}
	}
	w := cmd.cfg
	fmt.Fprintf(w, "{\n")

	// strings
	fmt.Fprintf(w, "  id: %#v,\n", key.GetHexKeyID())
	fmt.Fprintf(w, "  fingerprint: %#v,\n", key.GetFingerprint())

	// flags
	fmt.Fprintf(w, "  is_private: %#v,\n", key.IsPrivate())
	fmt.Fprintf(w, "  is_expired: %#v,\n", key.IsExpired())
	fmt.Fprintf(w, "  is_revoked: %#v,\n", key.IsRevoked())
	fmt.Fprintf(w, "  can_verify: %#v,\n", key.CanVerify())
	fmt.Fprintf(w, "  can_encrypt: %#v,\n", key.CanEncrypt())

	fmt.Fprintf(w, "}\n")
	return nil
}
