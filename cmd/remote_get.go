package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type RemoteGet struct {
	cfg   Config
	query string
}

func (cmd RemoteGet) Command() *cobra.Command {
	c := &cobra.Command{
		Use:     "get",
		Aliases: []string{"download", "pull", "g"},
		Args:    cobra.ExactArgs(1),
		Short:   "Download the key from a remote server",
		RunE: func(_ *cobra.Command, args []string) error {
			cmd.query = args[0]
			return cmd.run()
		},
	}
	return c
}

func (cmd RemoteGet) run() error {
	found := false
	keys := make(chan []byte)

	// run downloads from all servers
	servers := []Server{
		ServerGithub{},
		ServerKeybase{},
		ServerProtonmail{},
	}
	group := errgroup.Group{}
	runner := func(s Server) func() error {
		return func() error {
			key, err := s.Download(cmd.query)
			if key != nil {
				keys <- key
				found = true
			}
			return err
		}
	}
	for _, s := range servers {
		group.Go(runner(s))
	}

	// print all the keys that the servers returned
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for key := range keys {
			cmd.cfg.Write(key)
			cmd.cfg.Write([]byte{'\n'})
		}
		wg.Done()
	}()

	// wait for all servers and printer to finish
	err := group.Wait()
	close(keys)
	wg.Wait()
	if err != nil {
		return fmt.Errorf("cannot fetch the key: %v", err)
	}

	if !found {
		return errors.New("cannot find the key in any supported source")
	}
	return nil
}
