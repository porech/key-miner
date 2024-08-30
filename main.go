package main

import (
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func getPassFromTerminal() (string, error) {
	fmt.Fprintf(os.Stderr, "Enter passphrase: ")
	passphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintf(os.Stderr, "\n")
	if err != nil {
		return "", err
	}
	return string(passphrase), nil
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	username := os.Getenv("USER")
	if username == "" {
		username = "user"
	}
	defaultComment := username + "@" + hostname
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "This utility tries to find a ED25519 key pair whose public string has a specified sequence in it.\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <wanted-sequence>\n", os.Args[0])
		fmt.Printf("Example: %s LOVE\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Available flags:\n")
		flag.PrintDefaults()
	}

	keyPath := flag.String("path", "id_ed25519", "Path to write the generated private key (public key will be the same path with .pub extension)")
	keyComment := flag.String("comment", defaultComment, "Comment to add to the private and public key")
	emptyComment := flag.Bool("empty-comment", false, "Explicitly set an empty comment for the private and public key (overrides -comment)")
	keyPassphrase := flag.String("passphrase", "", "Passphrase to encrypt the private key (discouraged, use -ask-passphrase instead)")
	askPassphrase := flag.Bool("ask-passphrase", false, "Ask for the passphrase to encrypt the private key (overrides -passphrase)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: wanted-sequence is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Error: too many arguments\n")
		flag.Usage()
		os.Exit(1)
	}

	wantedSequence := flag.Arg(0)

	comment := *keyComment
	if *emptyComment {
		comment = ""
	}

	passphrase := *keyPassphrase
	if *askPassphrase {
		passphrase, err = getPassFromTerminal()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}
	publicKey, err := ssh.NewPublicKey(pub)
	if err != nil {
		panic(err)
	}
	pubkeyContent := base64.StdEncoding.EncodeToString(publicKey.Marshal())
	madeAttempts := 1
	for {
		if strings.Contains(pubkeyContent, wantedSequence) {
			break
		}
		pub, priv, err = ed25519.GenerateKey(nil)
		if err != nil {
			panic(err)
		}
		publicKey, err = ssh.NewPublicKey(pub)
		if err != nil {
			panic(err)
		}
		pubkeyContent = base64.StdEncoding.EncodeToString(publicKey.Marshal())
		madeAttempts++
		if madeAttempts % 100000 == 0 {
			fmt.Printf("%d attempts made\n", madeAttempts)
		}
	}
	var p *pem.Block
	if passphrase != "" {
		p, err = ssh.MarshalPrivateKeyWithPassphrase(crypto.PrivateKey(priv), comment, []byte(passphrase))
	} else {
		p, err = ssh.MarshalPrivateKey(crypto.PrivateKey(priv), comment)
	}
	if err != nil {
		panic(err)
	}
	privateKeyPem := pem.EncodeToMemory(p)
	publicKeyString := "ssh-ed25519" + " " + pubkeyContent
	if *keyComment != "" {
		publicKeyString += " " + comment
	}
	publicKeyString += "\n"
	os.WriteFile(*keyPath, privateKeyPem, 0600)
	os.WriteFile(*keyPath+".pub", []byte(publicKeyString), 0644)
	print("Private key written to: " + *keyPath + "\n")
}
