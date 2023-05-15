module github.com/bitweb-project/bted/bteutil/psbt

go 1.17

require (
	github.com/bitweb-project/bted v0.23.3
	github.com/bitweb-project/bted/btcec/v2 v2.1.3
	github.com/bitweb-project/bted/bteutil v1.1.0
	github.com/bitweb-project/bted/chaincfg/chainhash v1.0.2
	github.com/davecgh/go-spew v1.1.1
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/bitweb-project/yespower_go v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20200814200057-3d37ad5750ed // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/bitweb-project/bted/bteutil => ../

replace github.com/bitweb-project/bted => ../..
