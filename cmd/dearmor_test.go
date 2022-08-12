package cmd_test

import (
	"testing"

	"github.com/matryer/is"
)

func TestEncryptArmorDearmor(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	outE := CallHappy(t, "encrypt -p hi", []byte("message"))
	outA := CallHappy(t, "armor", outE)
	outD := CallHappy(t, "dearmor", outA)
	is.Equal(outD, outE)
}
