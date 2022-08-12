package cmd_test

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestKeyGenerate(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	outG := CallHappy(t, "key generate", nil)
	is.True(len(outG) > 4096)
}

func TestKeyGenerate_IsRandom(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	out1 := CallHappy(t, "key generate", nil)
	out2 := CallHappy(t, "key generate", nil)
	is.True(!bytes.Equal(out1, out2))
}
