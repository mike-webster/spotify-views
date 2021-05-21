package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mike-webster/spotify-views/encrypt"
	"github.com/mike-webster/spotify-views/keys"
)

func main() {
	mk, err := getMasterKey()
	if err != nil {
		panic(fmt.Sprint("couldnt find masterkey: ", err))
	}
	ctx := context.WithValue(context.Background(), keys.ContextMasterKey, mk)

	input, output, err := parseFilenames()
	if err != nil {
		panic(fmt.Sprint("couldnt read filenames: ", err))
	}

	err = act(ctx, input, output)
	if err != nil {
		panic(fmt.Sprint("couldnt act: ", err))
	}

	fmt.Println("fin")
}

func getMasterKey() (string, error) {
	mk := os.Getenv("MASTER_KEY")
	if len(mk) < 1 {
		return "", errors.New("no master key")
	}

	return mk, nil
}

func parseFilenames() (string, string, error) {
	if len(os.Args) < 4 {
		return "", "", errors.New("incorrect number of parameters")
	}
	args := os.Args[2:]

	if !strings.HasPrefix(args[0], "-in=") {
		return "", "", errors.New(fmt.Sprint("input flag invalid: ", args[0]))
	}

	if !strings.HasPrefix(args[1], "-out=") {
		return "", "", errors.New(fmt.Sprint("output flag invalid: ", args[1]))
	}

	input := strings.Replace(args[0], "-in=", "", 1)
	output := strings.Replace(args[1], "-out=", "", 1)
	return input, output, nil
}

func act(ctx context.Context, input, output string) error {
	if len(os.Args) < 2 {
		return errors.New("incorrect number of parameters")
	}

	direction := os.Args[1]

	if direction == "-e" {
		f, err := ioutil.ReadFile(input)
		if err != nil {
			panic(fmt.Sprint("couldn't read file: ", err))
		}

		b, err := encrypt.Encrypt(ctx, f)
		if err != nil {
			panic(fmt.Sprint("couldn't encrypt file: ", err))
		}

		err = ioutil.WriteFile(output, *b, 0644)
		if err != nil {
			panic(fmt.Sprint("couldn't write encrypted file: ", output, " - ", err))
		}

		fmt.Println("successfully encrypted and written to: ", output)
		return encryptFile(ctx, input, output)
	} else if direction == "-d" {
		return decryptFile(ctx, input, output)
	}

	panic(fmt.Sprint("unrecognized command: ", direction))
}

func encryptFile(ctx context.Context, input, output string) error {
	f, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	b, err := encrypt.Encrypt(ctx, f)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, *b, 0644)
	if err != nil {
		return err
	}

	fmt.Println("successfully encrypted and written to: ", output)
	return nil
}

func decryptFile(ctx context.Context, input, output string) error {
	f, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	b, err := encrypt.Decrypt(ctx, f)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(output, *b, 0644)
	if err != nil {
		return err
	}

	fmt.Println("successfully encrypted and written to: ", output)
	return nil
}
