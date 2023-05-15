package bteutil_test

import (
	"fmt"
	"math"

	"github.com/bitweb-project/bted/bteutil"
)

func ExampleAmount() {

	a := bteutil.Amount(0)
	fmt.Println("Zero Satoshi:", a)

	a = bteutil.Amount(1e8)
	fmt.Println("100,000,000 Satoshis:", a)

	a = bteutil.Amount(1e5)
	fmt.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Satoshi: 0 BTE
	// 100,000,000 Satoshis: 1 BTE
	// 100,000 Satoshis: 0.001 BTE
}

func ExampleNewAmount() {
	amountOne, err := bteutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := bteutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := bteutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := bteutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 BTE
	// 0.01234567 BTE
	// 0 BTE
	// invalid bitcoin amount
}

func ExampleAmount_unitConversions() {
	amount := bteutil.Amount(44433322211100)

	fmt.Println("Satoshi to kBTE:", amount.Format(bteutil.AmountKiloBTE))
	fmt.Println("Satoshi to BTE:", amount)
	fmt.Println("Satoshi to MilliBTE:", amount.Format(bteutil.AmountMilliBTE))
	fmt.Println("Satoshi to MicroBTE:", amount.Format(bteutil.AmountMicroBTE))
	fmt.Println("Satoshi to Satoshi:", amount.Format(bteutil.AmountSatoshi))

	// Output:
	// Satoshi to kBTE: 444.333222111 kBTE
	// Satoshi to BTE: 444333.222111 BTE
	// Satoshi to MilliBTE: 444333222.111 mBTE
	// Satoshi to MicroBTE: 444333222111 μBTE
	// Satoshi to Satoshi: 44433322211100 Satoshi
}
