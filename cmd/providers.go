package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Provider interface {
	Name() string
	Download(query string) ([]byte, error)
}

type ProviderGithub struct{}

func (ProviderGithub) Name() string {
	return "github"
}

func (ProviderGithub) Download(q string) ([]byte, error) {
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

type ProviderKeybase struct{}

func (ProviderKeybase) Name() string {
	return "keybase"
}

func (ProviderKeybase) Download(q string) ([]byte, error) {
	url := fmt.Sprintf("https://keybase.io/%s/pgp_keys.asc", q)
	return readURL(url)
}

type ProviderProtonmail struct{}

func (ProviderProtonmail) Name() string {
	return "protonmail"
}

func (ProviderProtonmail) Download(q string) ([]byte, error) {
	url := fmt.Sprintf("https://api.protonmail.ch/pks/lookup?op=get&search=%s", q)
	return readURL(url)
}

func readURL(url string) ([]byte, error) {
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
