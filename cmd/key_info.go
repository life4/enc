package cmd

import (
	"encoding/json"
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
		Example: "cat private.key | enc key info",
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
	ident := key.GetEntity().PrimaryIdentity()

	expirationStr := ""
	ttl := ident.SelfSignature.KeyLifetimeSecs
	if ttl != nil && *ttl != 0 {
		expiration := prim.CreationTime.Add(time.Duration(*ttl) * time.Second)
		expirationStr = expiration.Format(time.RFC3339)
	}
	result := map[string]interface{}{
		// basic key info
		"id":           key.GetHexKeyID(),
		"fingerprint":  key.GetFingerprint(),
		"algorithm":    cmd.algorithm(key),
		"created_at":   prim.CreationTime.Format(time.RFC3339),
		"expires_at":   expirationStr,
		"fingerprints": key.GetSHA256Fingerprints(),

		// user identity
		"identity": map[string]string{
			"name":    ident.UserId.Name,
			"email":   ident.UserId.Email,
			"comment": ident.UserId.Comment,
		},

		// flags
		"is_private":  key.IsPrivate(),
		"is_expired":  key.IsExpired(),
		"is_revoked":  key.IsRevoked(),
		"can_verify":  key.CanVerify(),
		"can_encrypt": key.CanEncrypt(),
	}

	b, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return fmt.Errorf("serialize JSON: %v", err)
	}
	_, err = cmd.cfg.Write(b)
	return err
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
