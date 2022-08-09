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
	c.AddCommand(Version{Stdout: w}.Command())
	c.AddCommand(Encrypt{Stdout: w, Stdin: r}.Command())
	c.AddCommand(Decrypt{Stdout: w, Stdin: r}.Command())
	c.AddCommand(Armor{Stdout: w, Stdin: r}.Command())
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
