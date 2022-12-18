package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type RemoteGet struct {
	cfg          Config
	query        string
	provider     string
	gitlabDomain string
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
	c.Flags().StringVarP(
		&cmd.provider, "provider", "p", "",
		"provider to use to get the key. All providers will be used if not specified",
	)
	c.Flags().StringVar(
		&cmd.gitlabDomain, "gitlab-domain", "gitlab.com",
		"domain name of GitLab server to use with gitlab provider",
	)
	return c
}

func (cmd RemoteGet) run() error {
	found := false
	keys := make(chan []byte)

	// run downloads from all providers
	group := errgroup.Group{}
	runner := func(p Provider) func() error {
		return func() error {
			key, err := p.Download(cmd.query)
			if key != nil {
				keys <- key
				found = true
			}
			return err
		}
	}
	providers := cmd.providers()
	if providers == nil {
		return fmt.Errorf("unknown provider: %s", cmd.provider)
	}
	for _, p := range providers {
		group.Go(runner(p))
	}

	// print all the keys that the providers returned
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for key := range keys {
			cmd.cfg.Write(key)
			cmd.cfg.Write([]byte{'\n'})
		}
		wg.Done()
	}()

	// wait for all providers and printer to finish
	err := group.Wait()
	close(keys)
	wg.Wait()
	if err != nil {
		return fmt.Errorf("cannot fetch the key: %v", err)
	}

	if !found {
		if cmd.provider == "" {
			return errors.New("cannot find the key in any supported source")
		} else {
			return errors.New("cannot find the key")
		}
	}
	return nil
}

func (cmd RemoteGet) providers() []Provider {
	providers := []Provider{
		ProviderGithub{},
		ProviderKeybase{},
		ProviderProtonmail{},
		ProviderGitlab{Domain: cmd.gitlabDomain},
	}
	if cmd.provider == "" {
		return providers
	}
	for _, p := range providers {
		if p.Name() == cmd.provider {
			return []Provider{p}
		}
	}
	return nil
}
