package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndValidateProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	err := os.WriteFile(path, []byte(`[profiles.default]
url = "https://netbird.example.test"
account_id = "account-1"
credential_ref = "env:NB_TEST_TOKEN"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	file, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	profile, err := file.Profile("default")
	if err != nil {
		t.Fatal(err)
	}
	if profile.AccountID != "account-1" || profile.CredentialRef != "env:NB_TEST_TOKEN" {
		t.Fatalf("unexpected profile: %+v", profile)
	}
}

func TestProfileReadOnlyFieldRoundTrips(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(`[profiles.production]
url = "https://netbird.example.test"
read_only = true
`), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	profile, err := file.Profile("production")
	if err != nil {
		t.Fatal(err)
	}
	if !profile.ReadOnly {
		t.Fatal("expected read-only profile")
	}
}

func TestDefaultStatePathUsesNBState(t *testing.T) {
	t.Setenv("NB_STATE", filepath.Join(t.TempDir(), "ledger.db"))
	if got, want := DefaultStatePath(), os.Getenv("NB_STATE"); got != want {
		t.Fatalf("state path=%q, want %q", got, want)
	}
}

func TestProfileRejectsInsecureRemoteURL(t *testing.T) {
	if err := (Profile{URL: "http://netbird.example.test"}).Validate(); err == nil {
		t.Fatal("expected insecure remote URL to fail")
	}
}

func TestResolveCredentialFromEnvironment(t *testing.T) {
	t.Setenv("NB_TEST_TOKEN", "secret-value")
	value, err := ResolveCredential("env:NB_TEST_TOKEN")
	if err != nil || value != "secret-value" {
		t.Fatalf("value=%q err=%v", value, err)
	}
}
