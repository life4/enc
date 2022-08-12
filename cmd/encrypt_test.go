package cmd_test

import (
	"testing"

	"github.com/matryer/is"
)

func TestEncryptPassword(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	outE := CallHappy(t, "encrypt -p hi", []byte("message"))
	is.Equal(len(outE), 106)
}
