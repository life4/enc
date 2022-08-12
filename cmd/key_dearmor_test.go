package cmd_test

import (
	"testing"

	"github.com/matryer/is"
)

func TestKeyGenerateArmorDearmor(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	outG := CallHappy(t, "key generate", nil)
	outA := CallHappy(t, "key armor", outG)
	outD := CallHappy(t, "key dearmor", outA)
	is.Equal(outG, outD)
}
