package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type Provider func(q string) ([]byte, error)

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

	// run all providers
	providers := []Provider{
		cmd.readGithub,
		cmd.readKeybase,
		cmd.readProtonmail,
	}
	group := errgroup.Group{}
	runner := func(p Provider) func() error {
		return func() error {
			key, err := p(cmd.query)
			if key != nil {
				keys <- key
				found = true
			}
			return err
		}
	}
	for _, p := range providers {
		group.Go(runner(p))
	}

	// print all the keys that the providers have found
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
		return errors.New("cannot find the key in any supported source")
	}
	return nil
}

func (RemoteGet) readGithub(q string) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/gpg_keys", q)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %v", err)
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("read response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	var keys []struct {
		Key string `json:"raw_key"`
	}
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return nil, fmt.Errorf("parse response: %v", err)
	}
	var buf bytes.Buffer
	for _, key := range keys {
		if key.Key != "" {
			buf.WriteString(key.Key)
		}
	}
	return io.ReadAll(&buf)
}

func (cmd RemoteGet) readKeybase(q string) ([]byte, error) {
	url := fmt.Sprintf("https://keybase.io/%s/pgp_keys.asc", q)
	return cmd.readURL(url)
}

func (cmd RemoteGet) readProtonmail(q string) ([]byte, error) {
	url := fmt.Sprintf("https://api.protonmail.ch/pks/lookup?op=get&search=%s", q)
	return cmd.readURL(url)
}

func (RemoteGet) readURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("read response: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
