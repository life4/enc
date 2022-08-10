package cmd

import "io"

// Config holds global configurations for commands
type Config struct {
	Stdin  io.Reader
	Stdout io.Writer
}

func (c Config) Read(p []byte) (n int, err error) {
	return c.Stdin.Read(p)
}

func (c Config) Write(p []byte) (n int, err error) {
	return c.Stdout.Write(p)
}
