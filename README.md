# key-miner

This utility tries to find a ED25519 key pair whose public string has a specified sequence in it.

## Why?

Good question. Vanity, probably? I happen to share my SSH public key with customers to let me access their servers, and I wanted it to contain a part of my name.

## Usage

Clone the repository and compile with:

```
go build
```

The only required argument is the string you want the public key to contain.

Example:

```
./key-miner LOVE
```

This will find a key pair whose public part contains the string "LOVE", and saves it to id_ed25519 and id_ed25519.pub.

Those are the optional flags:

```
  -ask-passphrase
    	Ask for the passphrase to encrypt the private key (overrides -passphrase)
  -comment string
    	Comment to add to the private and public key (default "user@hostname")
  -empty-comment
    	Explicitly set an empty comment for the private and public key (overrides -comment)
  -passphrase string
    	Passphrase to encrypt the private key (discouraged, use -ask-passphrase instead)
  -path string
    	Path to write the generated private key (public key will be the same path with .pub extension) (default "id_ed25519")
```

## Performance

It really depends on your CPU and, well, randomness. The correct key may pop up immediately or never, you can not know.
However, you should expect it to be really fast (less than one second) with up to 3 characters, and reasonably fast (some seconds) with 4 characters. I had success with 5 characters on a fast computer within 10 to 20 minutes. More than 5 characters are unlikely to produce a result in a reasonable amount of time, but you can try!
