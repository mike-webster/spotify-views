package env

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mike-webster/spotify-views/encrypt"
	"gopkg.in/yaml.v2"
)

type Secrets struct {
	LyricsKey    string `yaml:"lyrics_key"`
	DBHost       string `yaml:"db_host"`
	DBUser       string `yaml:"db_user"`
	DBPass       string `yaml:"db_pass"`
	DBName       string `yaml:"db_name"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	SecurityKey  string `yaml:"security_key"`
}

func (s *Secrets) IsValid() error {
	if len(s.LyricsKey) < 1 {
		return errors.New("missing lyrics key")
	}
	if len(s.DBHost) < 1 {
		return errors.New("missing db host")
	}
	if len(s.DBUser) < 1 {
		return errors.New("missing db user")
	}
	if len(s.DBPass) < 1 {
		return errors.New("missing db pass")
	}
	if len(s.DBName) < 1 {
		return errors.New("missing db name")
	}
	if len(s.ClientID) < 1 {
		return errors.New("missing client id")
	}
	if len(s.ClientSecret) < 1 {
		return errors.New("missing client secret")
	}

	return nil
}

var (
	secretsFile = "secrets.enc"
)

func ParseSecrets(ctx context.Context) (*Secrets, error) {
	env := os.Getenv("GO_ENV")
	if len(env) < 1 {
		return nil, errors.New("no go environment provided")
	}

	if !(env == "development" || env == "test" || env == "production" || env == "uat") {
		return nil, errors.New(fmt.Sprint("unsupported GO_ENV :", env))
	}

	var envs map[string]Secrets

	// Read the .enc file
	f, err := ioutil.ReadFile(secretsFile)
	if err != nil {
		return nil, err
	}

	// decrypt the bytes
	b, err := encrypt.Decrypt(ctx, f)
	if err != nil {
		return nil, err
	}

	// unmarshal the bytes
	err = yaml.Unmarshal(*b, &envs)
	if err != nil {
		return nil, err
	}

	// return the correct secrets
	ret := envs[env]

	err = ret.IsValid()
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
