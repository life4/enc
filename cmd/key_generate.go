package cmd

import (
	stdcrypto "crypto"
	"errors"
	"fmt"
	"os/user"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

type KeyGenerate struct {
	cfg     Config
	name    string
	email   string
	ktype   string
	comment string
	bits    int
	ttl     time.Duration
}

func (cmd KeyGenerate) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"create", "g"},
		Args:    cobra.NoArgs,
		Short:   "Generate new private key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.name, "name", "n", "", "your full name")
	c.Flags().StringVarP(&cmd.email, "email", "e", "", "your email address")
	c.Flags().StringVarP(&cmd.ktype, "type", "t", "rsa", "type of the key")
	c.Flags().IntVarP(&cmd.bits, "bits", "b", 4096, "size of RSA key in bits")
	c.Flags().DurationVar(&cmd.ttl, "ttl", 0, "validity period of the key")
	return c
}

func (cmd KeyGenerate) run() error {
	username := cmd.Username()
	if username == "" {
		return errors.New("--name is required")
	}
	alg := cmd.algorithm()
	if alg == 0 {
		return fmt.Errorf("unsupported key type: %v", cmd.ktype)
	}
	if alg != packet.PubKeyAlgoRSA && cmd.bits != 4096 {
		return errors.New("--bits has effect only with 'rsa' key type")
	}
	cfg := &packet.Config{
		Algorithm:              alg,
		RSABits:                cmd.bits,
		KeyLifetimeSecs:        uint32(cmd.ttl.Seconds()),
		Time:                   crypto.GetTime,
		DefaultHash:            stdcrypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
	}

	entity, err := openpgp.NewEntity(username, cmd.comment, cmd.email, cfg)
	if err != nil {
		return fmt.Errorf("cannot create entity: %v", err)
	}
	if entity.PrivateKey == nil {
		return errors.New("cannot generate private key")
	}
	err = entity.SerializePrivateWithoutSigning(cmd.cfg, nil)
	if err != nil {
		return fmt.Errorf("cannot serialize key: %v", err)
	}
	return nil
}

func (cmd KeyGenerate) Username() string {
	if cmd.name != "" {
		return cmd.name
	}
	user, err := user.Current()
	if err == nil {
		if user.Name != "" {
			return user.Name
		}
		return user.Username
	}
	return ""
}

func (cmd KeyGenerate) algorithm() packet.PublicKeyAlgorithm {
	switch strings.ToLower(cmd.ktype) {
	case "rsa":
		return packet.PubKeyAlgoRSA
	case "x25519", "eddsa":
		return packet.PubKeyAlgoEdDSA
	case "elgamal":
		return packet.PubKeyAlgoElGamal
	case "dsa":
		return packet.PubKeyAlgoDSA
	case "ecdh":
		return packet.PubKeyAlgoECDH
	case "ecdsa":
		return packet.PubKeyAlgoECDSA
	default:
		return 0
	}
}
