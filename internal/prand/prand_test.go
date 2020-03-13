package prand

import (
	"fmt"
	"math"
	"runtime"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func coinFlipStatus(testSeconds int, testFunc func() bool) (pctHeads float64) {
	var heads, tails int

	timer := time.NewTimer(time.Duration(testSeconds) * time.Second)

	threads := runtime.NumCPU() - 1
	fmt.Printf("Test will run for %d seconds using %d threads...\n", testSeconds, threads)

	//stop := make(chan bool, 1)

	for i := 0; i < threads; i++ {
		go func() {
			for {
				if testFunc() {
					heads++
					continue
				}
				tails++
			}
		}()
	}

	time.Sleep(time.Duration(2) * time.Second)
	fstr := "Heads: %d	Tails: %d	%% Heads: %.3f\n"
	for {
		select {
		case <-timer.C:
			return
		default:
			time.Sleep(time.Duration(500) * time.Millisecond)
			pctHeads = float64(heads) / float64(heads+tails)
			fmt.Printf(fstr, heads, tails, pctHeads*float64(100))

		}
	}
	return
}

func checkDrift(bias, winPct, maxDrift float64) (err error) {
	biasDec := bias / float64(100)
	drift := math.Abs(winPct - biasDec)

	fmt.Printf("\nDrift relative to %.2f was %.4f\n", biasDec, drift)
	fmt.Printf("Max allowed (%.4f)\n\n", maxDrift)

	if drift > maxDrift {
		err = errors.Errorf("Drift greater than max allowed (%.4f)", maxDrift)
	}
	return

}

func TestStringN(t *testing.T) {
	for i := 0; i < 100; i++ {
		str := StringN(i)
		if len(str) != i {
			t.Errorf("Expected string length %d, received %d", i, len(str))
		}

	}
	fmt.Println("Test StringN successful")
}

func TestCoinFlip(t *testing.T) {
	bias := float64(50)
	tf := func() bool {
		return CoinFlip()
	}
	winPct := coinFlipStatus(20, tf)
	err := checkDrift(bias, winPct, 0.007)
	if err != nil {
		t.Error(err)
	}

}

func TestCoinFlipWithBias(t *testing.T) {
	bias := float64(45)
	// CoinFlip is a convenience function for CoinFlipBias(50)
	tf := func() bool {
		return CoinFlipBias(bias)
	}
	winPct := coinFlipStatus(20, tf)
	err := checkDrift(bias, winPct, 0.007)
	if err != nil {
		t.Error(err)
	}

}
