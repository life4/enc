package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
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
	key, err := cmd.readGithub(cmd.query)
	if err != nil {
		return fmt.Errorf("cannot read key from github: %v", err)
	}
	if key != nil {
		cmd.cfg.Write(key)
		return nil
	}

	key, err = cmd.readKeybase(cmd.query)
	if err != nil {
		return fmt.Errorf("cannot read key from keybase: %v", err)
	}
	if key != nil {
		cmd.cfg.Write(key)
		return nil
	}

	key, err = cmd.readProtonmail(cmd.query)
	if err != nil {
		return fmt.Errorf("cannot read key from protonmail: %v", err)
	}
	if key != nil {
		cmd.cfg.Write(key)
		return nil
	}

	return errors.New("cannot find the key in any supported source")
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
	for _, key := range keys {
		if key.Key != "" {
			return []byte(key.Key), nil
		}
	}
	return nil, nil
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
