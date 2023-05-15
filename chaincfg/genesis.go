// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"time"

	"github.com/bitweb-project/bted/chaincfg/chainhash"
	"github.com/bitweb-project/bted/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x26, /* |.......&| */
				0x62, 0x69, 0x74, 0x77, 0x65, 0x62, 0x33, 0x2e, /* |bitweb3.| */
				0x30, 0x20, 0x3d, 0x20, 0x42, 0x54, 0x45, 0x2d, /* |0 = BTE-| */
				0x2d, 0x79, 0x65, 0x73, 0x70, 0x6f, 0x77, 0x65, /* |-yespowe| */
				0x72, 0x2d, 0x2d, 0x2d, 0x32, 0x30, 0x32, 0x31, /* |r---2021| */
				0x2f, 0x30, 0x35, 0x2f, 0x30, 0x33, /* |/05/03| */
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{
				0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
				0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
				0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
				0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
				0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
				0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
				0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
				0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
				0x1d, 0x5f, 0xac, /* |._.| */
			},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xf3, 0xf1, 0x62, 0xbe, 0x69, 0xa8, 0xf2, 0x2b,
	0x5a, 0x48, 0x4e, 0xf8, 0xd7, 0x5d, 0x43, 0xda,
	0xb9, 0xa7, 0x83, 0x7a, 0x76, 0xc3, 0x66, 0xd2,
	0x9b, 0x4d, 0xc5, 0x6b, 0x9c, 0x3e, 0x04, 0x00,
})

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xe0, 0xaa, 0xc3, 0x44, 0xbf, 0x21, 0x8b, 0xbb,
	0x5e, 0xee, 0x3d, 0x1d, 0x5c, 0xa0, 0x07, 0x57,
	0xa7, 0x86, 0x59, 0x85, 0x83, 0x78, 0x1c, 0x60,
	0xc5, 0x5e, 0x8f, 0xf0, 0x30, 0x33, 0x36, 0x04,

})

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1619971700, 0), // 2009-01-03 18:15:05 +0000 UTC
		Bits:       0x1f1fffff,               // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      651,               // 2083236893
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var regTestGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x32, 0x7f, 0x1b, 0x1d, 0x97, 0x6d, 0x22, 0x22,
	0x69, 0xe6, 0xb7, 0x48, 0x6c, 0x8f, 0xfc, 0x04,
	0x0e, 0x90, 0x19, 0xfd, 0x5d, 0x75, 0xd4, 0x81,
	0x74, 0x0e, 0xcb, 0x77, 0x7c, 0x1d, 0x78, 0x66,
})

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the regression test network.  It is the same as the merkle root for
// the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1619971818, 0), // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      1,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet3GenesisHash is the hash of the first block in the block chain for the
// test network (version 3).
var testNet3GenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x63, 0x4f, 0x18, 0x2f, 0x79, 0x5c, 0xb7, 0xaa,
	0x4a, 0x5c, 0x4f, 0x4b, 0x87, 0xe5, 0x65, 0xd9,
	0x2a, 0x53, 0x37, 0xfc, 0xe5, 0x56, 0x62, 0x5e,
	0xf9, 0xe4, 0x15, 0x2f, 0x54, 0x2a, 0x00, 0x00,
})

// testNet3GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 3).  It is the same as the merkle root
// for the main network.
var testNet3GenesisMerkleRoot = genesisMerkleRoot

// testNet3GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 3).
var testNet3GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},          // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNet3GenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1619971765, 0),  // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x1e3fffff,                // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      18156,                // 414098458
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// simNetGenesisHash is the hash of the first block in the block chain for the
// simulation test network.
var simNetGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x32, 0x7f, 0x1b, 0x1d, 0x97, 0x6d, 0x22, 0x22,
	0x69, 0xe6, 0xb7, 0x48, 0x6c, 0x8f, 0xfc, 0x04,
	0x0e, 0x90, 0x19, 0xfd, 0x5d, 0x75, 0xd4, 0x81,
	0x74, 0x0e, 0xcb, 0x77, 0x7c, 0x1d, 0x78, 0x66,
})

// simNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the simulation test network.  It is the same as the merkle root for
// the main network.
var simNetGenesisMerkleRoot = genesisMerkleRoot

// simNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the simulation test network.
var simNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: simNetGenesisMerkleRoot,  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1619971818, 0), // 2014-05-28 15:52:37 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      1,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// sigNetGenesisHash is the hash of the first block in the block chain for the
// signet test network.
var sigNetGenesisHash = chainhash.Hash{
	0x63, 0x4f, 0x18, 0x2f, 0x79, 0x5c, 0xb7, 0xaa,
	0x4a, 0x5c, 0x4f, 0x4b, 0x87, 0xe5, 0x65, 0xd9,
	0x2a, 0x53, 0x37, 0xfc, 0xe5, 0x56, 0x62, 0x5e,
	0xf9, 0xe4, 0x15, 0x2f, 0x54, 0x2a, 0x00, 0x00,
}

// sigNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the signet test network. It is the same as the merkle root for
// the main network.
var sigNetGenesisMerkleRoot = genesisMerkleRoot

// sigNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the signet test network.
var sigNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: sigNetGenesisMerkleRoot,  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1619971765, 0), // 2020-09-01 00:00:00 +0000 UTC
		Bits:       0x1e3fffff,               // 503543726 [00000377ae000000000000000000000000000000000000000000000000000000]
		Nonce:      18156,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}
