package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func Command(cfg Config) *cobra.Command {
	root := &cobra.Command{
		Use:   "enc",
		Short: "Enc is PGP for humans",
		Long: `
			A user-friendly CLI tool to work with PGP keys:
			create, add, list, encrypt, decrypt, sign, verify signatures.
		`,
		SilenceUsage: true, // do not print usage when command return an error
	}
	// $ enc version
	root.AddCommand(Version{cfg: cfg}.Command())
	// $ enc encrypt
	root.AddCommand(Encrypt{cfg: cfg}.Command())
	// $ enc decrypt
	root.AddCommand(Decrypt{cfg: cfg}.Command())
	// $ enc armor
	root.AddCommand(Armor{cfg: cfg}.Command())
	// $ enc dearmor
	root.AddCommand(Dearmor{cfg: cfg}.Command())

	sig := &cobra.Command{
		Use:     "sig",
		Aliases: []string{"s"},
		Short:   "Operations with signatures",
	}
	// $ enc sig create
	sig.AddCommand(SigCreate{cfg: cfg}.Command())
	// $ enc sig verify
	sig.AddCommand(SigVerify{cfg: cfg}.Command())
	// $ enc sig id
	sig.AddCommand(SigID{cfg: cfg}.Command())
	// $ enc sig armor
	sig.AddCommand(SigArmor{cfg: cfg}.Command())
	root.AddCommand(sig)

	key := &cobra.Command{
		Use:     "key",
		Aliases: []string{"k"},
		Short:   "Operations with key",
	}
	// $ enc key generate
	key.AddCommand(KeyGenerate{cfg: cfg}.Command())
	// $ enc key info
	key.AddCommand(KeyInfo{cfg: cfg}.Command())
	// $ enc key public
	key.AddCommand(KeyPublic{cfg: cfg}.Command())
	// $ enc key armor
	key.AddCommand(KeyArmor{cfg: cfg}.Command())
	// $ enc key dearmor
	key.AddCommand(KeyDearmor{cfg: cfg}.Command())
	// $ enc key lock
	key.AddCommand(KeyLock{cfg: cfg}.Command())
	// $ enc key unlock
	key.AddCommand(KeyUnlock{cfg: cfg}.Command())
	// $ enc key upload
	// ...
	// $ enc key download
	// ...
	// $ enc key revoke
	key.AddCommand(KeyRevoke{cfg: cfg}.Command())
	root.AddCommand(key)

	keys := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"keychain", "keyring", "c"},
		Short:   "Operations with key ring",
	}
	// $ enc keys list
	keys.AddCommand(KeysList{cfg: cfg}.Command())
	// $ enc keys get
	keys.AddCommand(KeysGet{cfg: cfg}.Command())
	// $ enc keys add
	// ...
	// $ enc keys delete
	// ...
	root.AddCommand(keys)

	remote := &cobra.Command{
		Use:     "remote",
		Aliases: []string{"r", "server"},
		Short:   "Operations with remote key servers",
	}
	// $ enc remote get
	remote.AddCommand(RemoteGet{cfg: cfg}.Command())
	root.AddCommand(remote)

	return root
}

func Main(args []string, r io.Reader, w io.Writer) error {
	cfg := Config{Stdin: r, Stdout: w}
	c := Command(cfg)
	c.SetArgs(args)
	return c.Execute()
}

func Entrypoint() {
	err := Main(os.Args[1:], os.Stdin, os.Stdout)
	if err != nil {
		os.Exit(1)
	}
}
