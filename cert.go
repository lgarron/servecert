package main

import (
	"errors"
	"os"
	"os/exec"
	"path"
)

func dataDir() string {
	// We'd use `os.UserDataDir()`, but that doesn't exist (yet?).
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homeDir, ".data/servecert")
}

func dataDirDescendant(ancestorPath string) string {
	return path.Join(dataDir(), ancestorPath)
}

func pathExists(fullPath string) bool {
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func mkcert(domain string) {
	err := os.MkdirAll(dataDir(), 0750)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(dataDirDescendant("root"), 0750)
	if err != nil {
		panic(err)
	}
	os.Setenv("CAROOT", dataDirDescendant("root"))

	certDir := path.Join(dataDirDescendant("certs"), domain)
	err = os.MkdirAll(certDir, 0750)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("mkcert", "-install", domain)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = certDir

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
