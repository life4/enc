package cmd_test

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestKeyGenerateArmor(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	outG := CallHappy(t, "key generate", nil)
	outA := CallHappy(t, "key armor", outG)
	is.True(bytes.HasPrefix(outA, []byte("---")))
}
