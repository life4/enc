package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func Command(r io.Reader, w io.Writer) *cobra.Command {
	c := &cobra.Command{
		Use:   "enc",
		Short: "Enc is PGP for humans",
		Long: `
			A user-friendly CLI tool to work with PGP keys:
			create, add, list, encrypt, decrypt, sign, verify signatures.
		`,
	}
	// enc version
	c.AddCommand(Version{Stdout: w}.Command())
	// enc encrypt
	c.AddCommand(Encrypt{Stdout: w, Stdin: r}.Command())
	// enc decrypt
	c.AddCommand(Decrypt{Stdout: w, Stdin: r}.Command())
	// enc armor
	c.AddCommand(Armor{Stdout: w, Stdin: r}.Command())
	// enc dearmor
	c.AddCommand(Dearmor{Stdout: w, Stdin: r}.Command())

	// enc key generate
	// enc key armor
	// enc key dearmor
	// enc key lock --pass
	// enc key unlock --pass
	// enc key fingerprints
	// enc key send

	// enc keyring list
	// enc keyring import
	// enc keyring export public
	// enc keyring export private
	// enc keyring delete
	return c
}

func Main(args []string, r io.Reader, w io.Writer) error {
	c := Command(r, w)
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
