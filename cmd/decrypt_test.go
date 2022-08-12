package cmd_test

import (
	"testing"

	"github.com/matryer/is"
)

func TestEncryptDecryptPassword(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	outE := CallHappy(t, "encrypt -p hi", []byte("message"))
	outD := CallHappy(t, "decrypt -p hi", outE)
	is.Equal(outD, []byte("message"))
}
