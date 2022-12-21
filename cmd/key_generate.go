package cmd

import (
	"bytes"
	stdcrypto "crypto"
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
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
	ttl     string
}

const (
	day  = 24 * time.Hour
	year = 8766 * time.Hour // average year is 365.25 days
)

func (cmd KeyGenerate) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"create", "g"},
		Args:    cobra.NoArgs,
		Short:   "Generate new private key",
		Example: "enc key generate > private.key",
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.run()
		},
	}
	c.Flags().StringVarP(&cmd.name, "name", "n", "", "your full name")
	c.Flags().StringVarP(&cmd.email, "email", "e", "", "your email address")
	c.Flags().StringVarP(&cmd.comment, "comment", "c", "", "a note to add to the key")
	c.Flags().StringVarP(&cmd.ktype, "type", "t", "rsa", "type of the key")
	c.Flags().IntVarP(&cmd.bits, "bits", "b", 4096, "size of RSA key in bits")
	c.Flags().StringVar(
		&cmd.ttl, "ttl", "1y",
		"validity period of the key. Can be a date (2020-12-30) or duration (4y30d, 24h)",
	)
	return c
}

func (cmd KeyGenerate) run() error {
	username := cmd.Username()
	if username == "" {
		return errors.New("--name is required")
	}
	email := cmd.Email()
	if email == "" {
		return errors.New("--email is required")
	}
	if cmd.ttl == "" {
		return errors.New("--ttl is required")
	}
	ttl, err := ParseDuration(cmd.ttl)
	if err != nil {
		return fmt.Errorf("cannot parse --ttl: %v", err)
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
		KeyLifetimeSecs:        uint32(ttl.Seconds()),
		Time:                   crypto.GetTime,
		DefaultHash:            stdcrypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
	}

	entity, err := openpgp.NewEntity(username, cmd.comment, email, cfg)
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

func (cmd KeyGenerate) Email() string {
	if cmd.email != "" {
		return cmd.email
	}
	c := exec.Command("git", "config", "user.email")
	var stdout bytes.Buffer
	c.Stdout = &stdout
	err := c.Run()
	if err == nil {
		return strings.TrimSpace(stdout.String())
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

func ParseDuration(ttl string) (time.Duration, error) {
	t, err := time.Parse("2006-01-02", ttl)
	if err == nil {
		return t.Sub(time.Now()), nil
	}

	var shift time.Duration

	// parse years
	parts := strings.Split(ttl, "y")
	if len(parts) == 2 {
		years, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("parse year: %v", err)
		}
		shift += time.Duration(years) * year
		ttl = parts[1]
	}

	// parse days
	parts = strings.Split(ttl, "d")
	if len(parts) == 2 {
		days, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, fmt.Errorf("parse year: %v", err)
		}
		shift += time.Duration(days) * day
		ttl = parts[1]
	}

	if ttl == "" {
		return shift, nil
	}
	d, err := time.ParseDuration(ttl)
	if err != nil {
		return 0, err
	}
	return d + shift, nil
}
