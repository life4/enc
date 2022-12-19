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

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadKeys(cfg Config) (*crypto.KeyRing, error) {
	var r io.Reader
	if cfg.HasStdin() {
		r = cfg
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("find home dir: %v", err)
		}
		f, err := os.Open(home + "/.gnupg/pubring.gpg")
		if err == os.ErrNotExist {
			return nil, fmt.Errorf("read from stdin: %v", err)
		}
		if err != nil {
			return nil, fmt.Errorf("open secring.gpg: %v", err)
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
		err = keyring.AddKey(key)
		if err != nil {
			return nil, fmt.Errorf("add key into keyring: %v", err)
		}
	}
	return keyring, nil
}

func ReadKeyStdin(cfg Config) (*crypto.Key, error) {
	if !cfg.HasStdin() {
		return nil, errors.New("no key passed into stdin")
	}
	return ReadKeyStream(cfg)
}

func ReadKeyFile(path string) (*crypto.Key, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %v", err)
	}
	defer f.Close()
	return ReadKeyStream(f)
}

func ReadKeyStream(r io.Reader) (*crypto.Key, error) {
	data, err := io.ReadAll(r)
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

func ReadPlainMessageStdin(cfg Config) (*crypto.PlainMessage, error) {
	if !cfg.HasStdin() {
		return nil, errors.New("no file passed into stdin")
	}
	data, err := io.ReadAll(cfg)
	if err != nil {
		return nil, fmt.Errorf("read from stdin: %v", err)
	}
	return crypto.NewPlainMessage(data), nil
}

func ReadPGPMessageStdin(cfg Config) (*crypto.PGPMessage, error) {
	data, err := io.ReadAll(cfg)
	if err != nil {
		return nil, fmt.Errorf("read from stdin: %v", err)
	}
	if bytes.HasPrefix(data, []byte("-----BEGIN PGP MESSAGE-----")) {
		message, err := crypto.NewPGPMessageFromArmored(string(data))
		if err != nil {
			return nil, fmt.Errorf("unarmor: %v", err)
		}
		return message, nil
	} else {
		return crypto.NewPGPMessage(data), nil
	}
}

func ReadSigStdin(cfg Config) (*crypto.PGPSignature, error) {
	if !cfg.HasStdin() {
		return nil, errors.New("no signature passed into stdin")
	}
	data, err := io.ReadAll(cfg)
	if err != nil {
		return nil, fmt.Errorf("read from stdin: %v", err)
	}
	isArmored := bytes.HasPrefix(data, []byte("-----BEGIN PGP"))
	if isArmored {
		sig, err := crypto.NewPGPSignatureFromArmored(string(data))
		if err != nil {
			return nil, fmt.Errorf("unarmor: %v", err)
		}
		return sig, nil
	} else {
		sig := crypto.NewPGPSignature(data)
		return sig, nil
	}
}
