# Homebrew release of Multihash Tool

This is the **unofficial** "Homebrew" version of the [`multihash`](https://multiformats.io/multihash/) tool command, which can be easily installed in non-Golang environments.

Its purpose is to check multi-hash values without the need to use or install IPFS.

```bash
# macOS and Linux (AMD64, ARM64/M1)
brew install KEINOS/apps/multihash
```

Same as "[go-multihash](https://github.com/multiformats/go-multihash)", but with the following differences:

- Version display.
- Unit tested.
- `golangci-lint` checked.
- `golint` checked.

## Usage

```bash
multihash [options] [FILE]
```

```shellsession
$ echo 'Hello, world!' | multihash
QmcwkKyBLujMQitrGSLdtFTzEYSzA7VcfARhFHbe4hZJc4

$ echo 'Hello, world!' | multihash -algorithm sha2-256 -encoding base58
QmcwkKyBLujMQitrGSLdtFTzEYSzA7VcfARhFHbe4hZJc4

$ echo 'Hello, world!' | multihash -a sha3-512 -e base64
FECMLmMwD5Ykttd2lf9/YCAcojWVCWxApTWrl425lyBO7BBmw/PULIaJWLu9+36c46LYg+GVEqkNlNvMksELCmQv
```

```shellsession
$ cat ./sample.txt
Hello, world!

$ multihash ./sample.txt
QmcwkKyBLujMQitrGSLdtFTzEYSzA7VcfARhFHbe4hZJc4
```

```shellsession
$ multihash -h
usage: multihash [options] [FILE]
Print or check multihash checksums.
With no FILE, or when FILE is -, read standard input.

Options:
  -a string
        one of: blake2b-128, blake2b-224, blake2b-256, blake2b-384, blake2b-512, blake2s-256, blake3, dbl-sha2-256, identity, keccak-224, keccak-256, keccak-384, keccak-512, md5, murmur3-x64-64, poseidon-bls12_381-a2-fc1, sha1, sha2-256, sha2-256-trunc254-padded, sha2-512, sha3, sha3-224, sha3-256, sha3-384, sha3-512, shake-128, shake-256, x11 (shorthand) (default "sha2-256")
  -algorithm string
        one of: blake2b-128, blake2b-224, blake2b-256, blake2b-384, blake2b-512, blake2s-256, blake3, dbl-sha2-256, identity, keccak-224, keccak-256, keccak-384, keccak-512, md5, murmur3-x64-64, poseidon-bls12_381-a2-fc1, sha1, sha2-256, sha2-256-trunc254-padded, sha2-512, sha3, sha3-224, sha3-256, sha3-384, sha3-512, shake-128, shake-256, x11 (default "sha2-256")
  -c string
        check checksum matches (shorthand)
  -check string
        check checksum matches
  -e string
        one of: raw, hex, base58, base64 (shorthand) (default "base58")
  -encoding string
        one of: raw, hex, base58, base64 (default "base58")
  -h    display help message (shorthand)
  -help
        display help message
  -l int
        checksums length in bits (truncate). -1 is default (shorthand) (default -1)
  -length int
        checksums length in bits (truncate). -1 is default (default -1)
  -q    quiet output (no newline on checksum, no error text) (shorthand)
  -quiet
        quiet output (no newline on checksum, no error text)
  -v    display app version (shorthand)
  -version
        display app version
```
