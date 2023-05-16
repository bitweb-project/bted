module github.com/bitweb-project/bted/btcec/v2

go 1.17

require (
	github.com/bitweb-project/bted/chaincfg/chainhash v1.0.2
	github.com/davecgh/go-spew v1.1.1
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1
)

require github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect

retract (
    v2.1.3 // invalid version
)
