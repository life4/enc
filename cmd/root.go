package cmd

import (
	"fmt"
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
	}
	// enc version
	root.AddCommand(Version{cfg: cfg}.Command())
	// enc encrypt
	root.AddCommand(Encrypt{cfg: cfg}.Command())
	// enc decrypt
	root.AddCommand(Decrypt{cfg: cfg}.Command())
	// enc armor
	root.AddCommand(Armor{cfg: cfg}.Command())
	// enc dearmor
	root.AddCommand(Dearmor{cfg: cfg}.Command())

	key := &cobra.Command{
		Use:   "key",
		Short: "Operations with a key",
	}
	// enc key generate
	key.AddCommand(KeyGenerate{cfg: cfg}.Command())
	// enc key armor
	// enc key dearmor
	// enc key lock --pass
	// enc key unlock --pass
	// enc key fingerprints
	// enc key send
	root.AddCommand(key)

	// enc keyring list
	// enc keyring import
	// enc keyring export public
	// enc keyring export private
	// enc keyring delete
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
