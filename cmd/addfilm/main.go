package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Neffats/final-scenes/models"
)

func main() {
	url := flag.String("url", "", "Youtube clip url")
	flag.Parse()

	if *name == "" {
		fmt.Println("missing arg: -url")
	}

	cmdArg := fmt.Sprintf("youtube-dl %s", *name)
	cmd := exec.Command(cmdArg)

	err := cmd.Run()

	if err != nil {
		fmt.Println("failed to run youtube-dl command: %v", err)
		return
	}
}
