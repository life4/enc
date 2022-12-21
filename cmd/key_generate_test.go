package cmd_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/life4/enc/cmd"
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

func TestParseDuration(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		given    string
		expected time.Duration
	}{
		{"4h", 4 * time.Hour},
		{"5d", 5 * 24 * time.Hour},
		{"5d4h", (5*24 + 4) * time.Hour},
		{"2y", 2 * 8766 * time.Hour},
		{"4h20m", 4*time.Hour + 20*time.Minute},
	}
	for _, tCase := range testCases {
		tc := tCase
		t.Run(tc.given, func(t *testing.T) {
			is := is.New(t)
			t.Parallel()
			actual, err := cmd.ParseDuration(tc.given)
			is.NoErr(err)
			is.Equal(actual, tc.expected)
		})
	}
}
