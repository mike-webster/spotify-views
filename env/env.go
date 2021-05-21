package env

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Host string `envconfig:"HOST"`
	Port string `envconfig:"PORT"`
}

func (e *Env) IsValid() error {
	if len(e.Host) < 1 {
		return errors.New("no host provided")
	}

	if len(e.Port) < 1 {
		return errors.New("no port provided")
	}

	return nil
}

func ParseEnv() (*Env, error) {
	e := Env{}
	envconfig.MustProcess("", &e)
	err := e.IsValid()
	if err != nil {
		return nil, err
	}

	return &e, nil
}
