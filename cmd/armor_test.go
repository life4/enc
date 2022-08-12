package cmd_test

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestEncryptArmor(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	outE := CallHappy(t, "encrypt -p hi", []byte("message"))
	outA := CallHappy(t, "armor", outE)
	is.True(bytes.HasPrefix(outA, []byte("---")))
}
