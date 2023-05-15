// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"math"
	"math/big"
	"time"

	"github.com/bitweb-project/bted/chaincfg/chainhash"
)

var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// oneLsh256 is 1 shifted left 256 bits.  It is defined here to avoid
	// the overhead of creating it multiple times.
	oneLsh256 = new(big.Int).Lsh(bigOne, 256)
)

// HashToBig converts a chainhash.Hash into a big.Int that can be used to
// perform math comparisons.
func HashToBig(hash *chainhash.Hash) *big.Int {
	// A Hash is in little-endian, but the big package wants the bytes in
	// big-endian, so reverse them.
	buf := *hash
	blen := len(buf)
	for i := 0; i < blen/2; i++ {
		buf[i], buf[blen-1-i] = buf[blen-1-i], buf[i]
	}

	return new(big.Int).SetBytes(buf[:])
}

// CompactToBig converts a compact representation of a whole number N to an
// unsigned 32-bit number.  The representation is similar to IEEE754 floating
// point numbers.
//
// Like IEEE754 floating point, there are three basic components: the sign,
// the exponent, and the mantissa.  They are broken out as follows:
//
//	* the most significant 8 bits represent the unsigned base 256 exponent
// 	* bit 23 (the 24th bit) represents the sign bit
//	* the least significant 23 bits represent the mantissa
//
//	-------------------------------------------------
//	|   Exponent     |    Sign    |    Mantissa     |
//	-------------------------------------------------
//	| 8 bits [31-24] | 1 bit [23] | 23 bits [22-00] |
//	-------------------------------------------------
//
// The formula to calculate N is:
// 	N = (-1^sign) * mantissa * 256^(exponent-3)
//
// This compact form is only used in bitcoin to encode unsigned 256-bit numbers
// which represent difficulty targets, thus there really is not a need for a
// sign bit, but it is implemented here to stay consistent with bitcoind.
func CompactToBig(compact uint32) *big.Int {
	// Extract the mantissa, sign bit, and exponent.
	mantissa := compact & 0x007fffff
	isNegative := compact&0x00800000 != 0
	exponent := uint(compact >> 24)

	// Since the base for the exponent is 256, the exponent can be treated
	// as the number of bytes to represent the full 256-bit number.  So,
	// treat the exponent as the number of bytes and shift the mantissa
	// right or left accordingly.  This is equivalent to:
	// N = mantissa * 256^(exponent-3)
	var bn *big.Int
	if exponent <= 3 {
		mantissa >>= 8 * (3 - exponent)
		bn = big.NewInt(int64(mantissa))
	} else {
		bn = big.NewInt(int64(mantissa))
		bn.Lsh(bn, 8*(exponent-3))
	}

	// Make it negative if the sign bit is set.
	if isNegative {
		bn = bn.Neg(bn)
	}

	return bn
}

// BigToCompact converts a whole number N to a compact representation using
// an unsigned 32-bit number.  The compact representation only provides 23 bits
// of precision, so values larger than (2^23 - 1) only encode the most
// significant digits of the number.  See CompactToBig for details.
func BigToCompact(n *big.Int) uint32 {
	// No need to do any work if it's zero.
	if n.Sign() == 0 {
		return 0
	}

	// Since the base for the exponent is 256, the exponent can be treated
	// as the number of bytes.  So, shift the number right or left
	// accordingly.  This is equivalent to:
	// mantissa = mantissa / 256^(exponent-3)
	var mantissa uint32
	exponent := uint(len(n.Bytes()))
	if exponent <= 3 {
		mantissa = uint32(n.Bits()[0])
		mantissa <<= 8 * (3 - exponent)
	} else {
		// Use a copy to avoid modifying the caller's original number.
		tn := new(big.Int).Set(n)
		mantissa = uint32(tn.Rsh(tn, 8*(exponent-3)).Bits()[0])
	}

	// When the mantissa already has the sign bit set, the number is too
	// large to fit into the available 23-bits, so divide the number by 256
	// and increment the exponent accordingly.
	if mantissa&0x00800000 != 0 {
		mantissa >>= 8
		exponent++
	}

	// Pack the exponent, sign bit, and mantissa into an unsigned 32-bit
	// int and return it.
	compact := uint32(exponent<<24) | mantissa
	if n.Sign() < 0 {
		compact |= 0x00800000
	}
	return compact
}

// CalcWork calculates a work value from difficulty bits.  Bitcoin increases
// the difficulty for generating a block by decreasing the value which the
// generated hash must be less than.  This difficulty target is stored in each
// block header using a compact representation as described in the documentation
// for CompactToBig.  The main chain is selected by choosing the chain that has
// the most proof of work (highest difficulty).  Since a lower target difficulty
// value equates to higher actual difficulty, the work value which will be
// accumulated must be the inverse of the difficulty.  Also, in order to avoid
// potential division by zero and really small floating point numbers, the
// result adds 1 to the denominator and multiplies the numerator by 2^256.
func CalcWork(bits uint32) *big.Int {
	// Return a work value of zero if the passed difficulty bits represent
	// a negative number. Note this should not happen in practice with valid
	// blocks, but an invalid block could trigger it.
	difficultyNum := CompactToBig(bits)
	if difficultyNum.Sign() <= 0 {
		return big.NewInt(0)
	}

	// (1 << 256) / (difficultyNum + 1)
	denominator := new(big.Int).Add(difficultyNum, bigOne)
	return new(big.Int).Div(oneLsh256, denominator)
}

// calcEasiestDifficulty calculates the easiest possible difficulty that a block
// can have given starting difficulty bits and a duration.  It is mainly used to
// verify that claimed proof of work by a block is sane as compared to a
// known good checkpoint.
func (b *BlockChain) calcEasiestDifficulty(bits uint32, duration time.Duration) uint32 {
	// Convert types used in the calculations below.
	durationVal := int64(duration / time.Second)
	adjustmentFactor := big.NewInt(b.chainParams.RetargetAdjustmentFactor)

	// The test network rules allow minimum difficulty blocks after more
	// than twice the desired amount of time needed to generate a block has
	// elapsed.
	if b.chainParams.ReduceMinDifficulty {
		reductionTime := int64(b.chainParams.MinDiffReductionTime /
			time.Second)
		if durationVal > reductionTime {
			return b.chainParams.PowLimitBits
		}
	}

	// Since easier difficulty equates to higher numbers, the easiest
	// difficulty for a given duration is the largest value possible given
	// the number of retargets for the duration and starting difficulty
	// multiplied by the max adjustment factor.
	newTarget := CompactToBig(bits)
	for durationVal > 0 && newTarget.Cmp(b.chainParams.PowLimit) < 0 {
		newTarget.Mul(newTarget, adjustmentFactor)
		durationVal -= b.maxRetargetTimespan
	}

	// Limit new value to the proof of work limit.
	if newTarget.Cmp(b.chainParams.PowLimit) > 0 {
		newTarget.Set(b.chainParams.PowLimit)
	}

	return BigToCompact(newTarget)
}

// findPrevTestNetDifficulty returns the difficulty of the previous block which
// did not have the special testnet minimum difficulty rule applied.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) findPrevTestNetDifficulty(startNode *blockNode) uint32 {
	// Search backwards through the chain for the last block without
	// the special rule applied.
	iterNode := startNode
	for iterNode != nil && iterNode.height%b.blocksPerRetarget != 0 &&
		iterNode.bits == b.chainParams.PowLimitBits {

		iterNode = iterNode.parent
	}

	// Return the found difficulty or the minimum difficulty if no
	// appropriate block was found.
	lastBits := b.chainParams.PowLimitBits
	if iterNode != nil {
		lastBits = iterNode.bits
	}
	return lastBits
}

// calcNextRequiredDifficulty calculates the required difficulty for the block
// after the passed previous block node based on the difficulty retarget rules.
// This function differs from the exported CalcNextRequiredDifficulty in that
// the exported version uses the current best chain as the previous block node
// while this function accepts any block node.
func (b *BlockChain) calcNextRequiredDifficulty(lastNode *blockNode, newBlockTime time.Time) (uint32, error) {
	// Genesis block.
	if lastNode.height > 90 {
		return b.lwmaCalculateNextWorkRequired(lastNode, newBlockTime)
	}

	if (lastNode.height+1)%b.blocksPerRetarget != 0 {
		// Return the same difficulty
		return lastNode.bits, nil
	}

	// Go back by what we want to be b.blocksPerRetarget worth of blocks
	nHeightFirst := lastNode.height - (b.blocksPerRetarget - 1)
	if nHeightFirst < 0 {
		// If we don't have enough blocks, return the pow limit
		return b.chainParams.PowLimitBits, nil
	}

	// Get the ancestor block at nHeightFirst
	firstNode := lastNode.Ancestor(b.blocksPerRetarget - 1)
	if firstNode == nil {
		return 0, AssertError("unable to obtain previous retarget block")
	}

	// Check if the previous block has the expected difficulty
	if lastNode.bits != firstNode.bits {
		// If not, return the difficulty of the previous block
		return lastNode.bits, nil
	}

	// Calculate the actual and adjusted time spans
	firstTime := time.Unix(firstNode.timestamp, 0)
	lastTime := time.Unix(lastNode.timestamp, 0)
	actualTimespan := lastTime.Sub(firstTime)
	
	// Clamp the actual time passed within the retarget window to the max and
	// min retarget timespan values.
	adjustedTimespan := float64(actualTimespan)
	if adjustedTimespan < float64(b.minRetargetTimespan)*float64(time.Second) {
		adjustedTimespan = float64(b.minRetargetTimespan) * float64(time.Second)
	} else if adjustedTimespan > float64(b.maxRetargetTimespan)*float64(time.Second) {
		adjustedTimespan = float64(b.maxRetargetTimespan) * float64(time.Second)
	}
	
	// Calculate the new target difficulty value based on the previous
	// target and the amount of time elapsed in the current retarget
	// interval.
	oldTarget := CompactToBig(lastNode.bits)
	newTarget := new(big.Int).Mul(oldTarget, big.NewInt(int64(adjustedTimespan)))
	targetTimeSpan := int64(b.chainParams.TargetTimePerBlock.Seconds()) * int64(b.blocksPerRetarget)
	newTarget.Div(newTarget, big.NewInt(targetTimeSpan))
	
	if newTarget.Cmp(b.chainParams.PowLimit) > 0 {
		newTarget.Set(b.chainParams.PowLimit)
	}
	newTargetBits := BigToCompact(newTarget)
	
	// Log debug information about the new target.
	log.Debugf("Difficulty retarget at block height %d", lastNode.height+1)
	log.Debugf("Old target %08x (%064x)", lastNode.bits, oldTarget)
	log.Debugf("New target %08x (%064x)", newTargetBits, CompactToBig(newTargetBits))
	log.Debugf("Actual timespan %v, adjusted timespan %v, target timespan %v",
		actualTimespan,
		time.Duration(int64(adjustedTimespan)),
		b.chainParams.TargetTimespan)
	
	return newTargetBits, nil
}


func (b *BlockChain) lwmaCalculateNextWorkRequired(lastNode *blockNode, newBlockTime time.Time) (uint32, error) {
	const T int64 = 60
	const N int64 = 90
	const k int64 = N * (N + 1) * T / 2
	height := int64(lastNode.height)
	if height <= N {
		log.Debugf("lwmaCalculateNextWorkRequired: block height %d must be greater than %d", height, N)
		return 0, nil
	}

	sumTarget := big.NewInt(0)
	t := int64(0)
	j := int64(0)
	solvetime := int64(0)

	for i := height - N + 1; i <= height; i++ {
		blockNode := lastNode.Ancestor(int32(i))
		blockNodePrev := blockNode.Ancestor(int32(i - 1))

		solvetime = blockNode.timestamp - blockNodePrev.timestamp
		solvetime = int64(math.Max(float64(-6*T), math.Min(float64(solvetime), float64(6*T))))
		log.Debugf("lwmaCalculateNextWorkRequired: block height=%d, solvetime=%d", blockNode.height, solvetime)
		j++
		t += solvetime * j
		target := CompactToBig(blockNode.bits)
		sumTarget.Add(sumTarget, target.Div(target, big.NewInt(int64(k*N))))
		log.Debugf("lwmaCalculateNextWorkRequired: block height=%d, target=%d, sumTarget=%d", blockNode.height, target, sumTarget)

	}

	// Calculate the next target using the LWMA formula
	if t < k/10 {
		t = k/10
	}
	nextTarget := new(big.Int).Mul(big.NewInt(t), sumTarget)
	log.Debugf("lwmaCalculateNextWorkRequired: t=%d, sumTarget=%d, nextTarget=%d", t, sumTarget, nextTarget)

	// Set the next target to the maximum difficulty if it exceeds it
	if nextTarget.Cmp(b.chainParams.PowLimit) > 0 {
		nextTarget = b.chainParams.PowLimit
	}

	log.Debugf("lwmaCalculateNextWorkRequired: height=%d, target=%d", height+1, BigToCompact(nextTarget))
	return BigToCompact(nextTarget), nil
}


// CalcNextRequiredDifficulty calculates the required difficulty for the block
// after the end of the current best chain based on the difficulty retarget
// rules.
//
// This function is safe for concurrent access.
func (b *BlockChain) CalcNextRequiredDifficulty(timestamp time.Time) (uint32, error) {
	b.chainLock.Lock()
	difficulty, err := b.calcNextRequiredDifficulty(b.bestChain.Tip(), timestamp)
	b.chainLock.Unlock()
	return difficulty, err
}
