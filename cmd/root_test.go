package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/life4/enc/cmd"
	"github.com/matryer/is"
)

func CallHappy(t *testing.T, sargs string, stdin []byte) []byte {
	is := is.New(t)
	r := bytes.Buffer{}
	_, err := r.Write(stdin)
	is.NoErr(err)
	w := bytes.Buffer{}
	args := strings.Split(sargs, " ")
	err = cmd.Main(args, &r, &w)
	is.NoErr(err)
	return w.Bytes()
}
