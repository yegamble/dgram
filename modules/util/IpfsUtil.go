package util

import (
	ipfs "github.com/ipfs/go-ipfs-api"
	"log"
	"os"
)

func UploadToIPFS(path string) (string, error) {
	shell := ipfs.NewShell("localhost:5001")

	bufferFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	hash, err := shell.Add(bufferFile)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	shell.Pin(hash)

	return hash, nil
}
