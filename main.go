package main

import (
	"crypto/cipher"
	"fmt"
	"math/rand"
	"time"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/util/random"
)

func main() {
	rand.Seed(time.Now().Unix())
	//2 - fee amount, 3 - splitting amount in 3 amounts
	t := newTxWorker(int64(2), 3)
	//70 - transaction amount
	t.ProcessTxsBlock(70)
}

type txworker struct {
	//fee
	fee int64
	//number of small transactions (splitting amount in 3 amounts)
	txCount int64
	//curve 25519 library
	suite *edwards25519.SuiteEd25519
	//public curve points
	G, H kyber.Point
	//random number generated
	rng cipher.Stream
}

// constructor
func newTxWorker(fee int64, txCount int64) *txworker {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	rng := random.New()
	return &txworker{
		fee:     fee,
		txCount: txCount,
		suite:   suite,
		G:       suite.Point().Pick(rng),
		H:       suite.Point().Pick(rng),
		rng:     rng,
	}
}

func (w *txworker) ProcessTxsBlock(amount int64) {
	inputs := splitAmount(amount, w.txCount)
	outputs := splitAmount(amount-w.fee, w.txCount-1)
	outputs = append(outputs, w.fee)
	inputCommitment, inputR := w.sumTxs(inputs)
	outputCommitment, outputR := w.sumTxs(outputs[:len(outputs)-1])
	feeCommitment := w.calcFeeComm(inputR, outputR)
	fmt.Printf("Input commitment: %s\n", inputCommitment)
	fmt.Printf("Inputs : %v %v\n ", inputs, sumSlice(inputs))
	fmt.Printf("Output commitment: %s\n", w.suite.Point().Add(outputCommitment, feeCommitment))
	fmt.Printf("Outputs : %v %v\n ", outputs, sumSlice(outputs))
}

func (w *txworker) sumTxs(txs []int64) (kyber.Point, kyber.Scalar) {
	var sumCommitment kyber.Point
	// to find random for fee
	var sumRandom kyber.Scalar
	for i, value := range txs {
		r := w.suite.Scalar().Pick(w.suite.RandomStream())
		a := w.suite.Scalar().SetInt64(value)

		rG := w.suite.Point().Mul(r, w.G)
		aH := w.suite.Point().Mul(a, w.H)
		c := w.suite.Point().Add(rG, aH)

		if i == 0 {
			sumRandom = r.Clone()
			sumCommitment = c.Clone()
			continue
		}
		sumRandom = w.suite.Scalar().Add(sumRandom, r)
		sumCommitment = w.suite.Point().Add(sumCommitment, c)
	}

	return sumCommitment, sumRandom
}

// finding a fee commitment
func (w *txworker) calcFeeComm(in, out kyber.Scalar) kyber.Point {
	r := w.suite.Scalar().Sub(in, out)
	a := w.suite.Scalar().SetInt64(w.fee)
	rG := w.suite.Point().Mul(r, w.G)
	aH := w.suite.Point().Mul(a, w.H)
	comm := w.suite.Point().Add(rG, aH)

	return comm
}

// splitting the amount of transaction in 3 small amounts (for inputs and outputs)
func splitAmount(value int64, size int64) []int64 {
	values := make([]int64, 0, size)
	for i := int64(0); i < size-1; i++ {
		limit := value - size + i
		r := rand.Int63n(limit) + 1
		value -= r
		values = append(values, r)
	}
	values = append(values, value)
	return values
}

// for printing
func sumSlice(s []int64) int64 {
	var sum int64
	for _, v := range s {
		sum += v
	}
	return sum
}
