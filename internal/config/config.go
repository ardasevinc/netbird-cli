package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type File struct {
	Profiles map[string]Profile `toml:"profiles"`
}

type Profile struct {
	URL            string `toml:"url"`
	AccountID      string `toml:"account_id"`
	CredentialRef  string `toml:"credential_ref"`
	CAFile         string `toml:"ca_file"`
	ServerIdentity string `toml:"server_identity"`
	ReadOnly       bool   `toml:"read_only"`
}

func DefaultPath() string {
	if value := os.Getenv("NB_CONFIG"); value != "" {
		return value
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(".config", "nb", "config.toml")
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "nb", "config.toml")
}

func DefaultStatePath() string {
	if value := os.Getenv("NB_STATE"); value != "" {
		return value
	}
	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join(".local", "state", "nb", "ledger.db")
		}
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, "nb", "ledger.db")
}

func Load(path string) (File, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- the operator explicitly selects the local config path.
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return File{Profiles: map[string]Profile{}}, nil
		}
		return File{}, fmt.Errorf("read config %s: %w", path, err)
	}
	var file File
	if err := toml.Unmarshal(data, &file); err != nil {
		return File{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	if file.Profiles == nil {
		file.Profiles = map[string]Profile{}
	}
	return file, nil
}

func (f File) Profile(name string) (Profile, error) {
	profile, ok := f.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}
	if err := profile.Validate(); err != nil {
		return Profile{}, fmt.Errorf("profile %q: %w", name, err)
	}
	return profile, nil
}

func (p Profile) Validate() error {
	if p.URL == "" {
		return errors.New("url is required")
	}
	u, err := url.Parse(p.URL)
	if err != nil || u.Scheme == "" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("url must be an absolute origin without credentials, query, or fragment")
	}
	if u.Scheme != "https" && !isLoopback(u.Hostname()) {
		return errors.New("url must use https unless it targets loopback")
	}
	if p.CredentialRef != "" && !strings.HasPrefix(p.CredentialRef, "env:") && !strings.HasPrefix(p.CredentialRef, "file:") {
		return errors.New("credential_ref must use an env or file reference")
	}
	if strings.HasPrefix(p.CredentialRef, "env:") && len(strings.TrimPrefix(p.CredentialRef, "env:")) == 0 {
		return errors.New("credential_ref env name is empty")
	}
	if strings.HasPrefix(p.CredentialRef, "file:") && len(strings.TrimPrefix(p.CredentialRef, "file:")) == 0 {
		return errors.New("credential_ref file path is empty")
	}
	return nil
}

func isLoopback(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
