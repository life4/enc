package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func ReadKeys(cfg Config) (*crypto.KeyRing, error) {
	var r io.Reader
	if cfg.HasStdin() {
		r = cfg
	} else {
		home, err := os.UserHomeDir()
		f, err := os.Open(home + "/.gnupg/pubring.gpg")
		if err != nil {
			return nil, fmt.Errorf("find home dir: %v", err)
		}
		if err == os.ErrNotExist {
			return nil, fmt.Errorf("read from stdin: %v", err)
		}
		if err != nil {
			return nil, fmt.Errorf("read secring.gpg: %v", err)
		}
		r = f
		defer f.Close()
	}
	entities, err := openpgp.ReadKeyRing(r)
	if err != nil {
		return nil, fmt.Errorf("read keyring: %v", err)
	}
	keyring, _ := crypto.NewKeyRing(nil)
	for _, entity := range entities {
		key, err := crypto.NewKeyFromEntity(entity)
		if err != nil {
			return nil, fmt.Errorf("parse key: %v", err)
		}
		keyring.AddKey(key)
	}
	return keyring, nil
}

func ReadKey(cfg Config) (*crypto.Key, error) {
	if !cfg.HasStdin() {
		return nil, errors.New("no key passed into stdin")
	}
	data, err := io.ReadAll(cfg)
	if err != nil {
		return nil, fmt.Errorf("read from stdin: %v", err)
	}
	isArmored := bytes.HasPrefix(data, []byte("-----BEGIN PGP"))
	if isArmored {
		key, err := crypto.NewKeyFromArmored(string(data))
		if err != nil {
			return nil, fmt.Errorf("unarmor key: %v", err)
		}
		return key, nil
	} else {
		key, err := crypto.NewKey(data)
		if err != nil {
			return nil, fmt.Errorf("parse key: %v", err)
		}
		return key, nil
	}
}
