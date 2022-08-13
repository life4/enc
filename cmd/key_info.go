package cmd

import (
	"fmt"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type KeyInfo struct {
	cfg Config
}

func (cmd KeyInfo) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "info",
		Aliases: []string{"inspect", "i"},
		Args:    cobra.NoArgs,
		Short:   "Show information about key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	return c
}

func (cmd KeyInfo) run() error {
	key, err := ReadKeyStdin(cmd.cfg)
	if err != nil {
		return fmt.Errorf("cannot read key: %v", err)
	}
	prim := key.GetEntity().PrimaryKey
	w := cmd.cfg
	fmt.Fprintf(w, "{\n")

	// strings
	fmt.Fprintf(w, "  id: %#v,\n", key.GetHexKeyID())
	fmt.Fprintf(w, "  fingerprint: %#v,\n", key.GetFingerprint())
	fmt.Fprintf(w, "  algorithm: %#v,\n", cmd.algorithm(key))
	fmt.Fprintf(w, "  created_at: %#v,\n", prim.CreationTime.Format(time.RFC3339))

	// user identity
	ident := key.GetEntity().PrimaryIdentity()
	fmt.Fprintln(w, "  identity: {")
	fmt.Fprintf(w, "    name: %#v,\n", ident.UserId.Name)
	fmt.Fprintf(w, "    email: %#v,\n", ident.UserId.Email)
	fmt.Fprintf(w, "    comment: %#v\n", ident.UserId.Comment)
	fmt.Fprintln(w, "  },")

	// flags
	fmt.Fprintf(w, "  is_private: %#v,\n", key.IsPrivate())
	fmt.Fprintf(w, "  is_expired: %#v,\n", key.IsExpired())
	fmt.Fprintf(w, "  is_revoked: %#v,\n", key.IsRevoked())
	fmt.Fprintf(w, "  can_verify: %#v,\n", key.CanVerify())
	fmt.Fprintf(w, "  can_encrypt: %#v\n", key.CanEncrypt())

	fmt.Fprintf(w, "}\n")
	return nil
}

func (KeyInfo) algorithm(key *crypto.Key) string {
	switch key.GetEntity().PrimaryKey.PubKeyAlgo {
	case packet.PubKeyAlgoDSA:
		return "DSA"
	case packet.PubKeyAlgoECDH:
		return "ECDH"
	case packet.PubKeyAlgoECDSA:
		return "ECDSA"
	case packet.PubKeyAlgoEdDSA:
		return "EdDSA"
	case packet.PubKeyAlgoElGamal:
		return "ElGamal"
	case packet.PubKeyAlgoRSA:
		return "RSA"
	default:
		return "?"
	}
}
