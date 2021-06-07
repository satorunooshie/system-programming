package main

import (
	"fmt"
	"github.com/tmc/keyring"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func main() {
	secretValue, err := keyring.Get("progo-keyring-test", "password")
	if err == keyring.ErrNotFound {
		// 未登録
		fmt.Printf("Secret Value is not found. Please Type:")
		pw, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			panic(err)
		}
		// 登録
		if err := keyring.Set("progo-keyring-test", "password", string(pw)); err != nil {
			panic(err)
		}
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("Secret Value: %s\n", secretValue)
}
