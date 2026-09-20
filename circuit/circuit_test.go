package circuit_test

import (
	"testing"

	"github.com/DylanSp/go-logical-circuits/circuit"
	"github.com/DylanSp/go-logical-circuits/circuit/wire"
	"github.com/stretchr/testify/assert"
)

func TestBasicGates(t *testing.T) {
	t.Run("NOT gate", func(t *testing.T) {
		// set up gate
		cir := circuit.NewCircuit()
		inWire := cir.AddInputWire("input")
		outWire := cir.Not(inWire)

		cir.Initialize()

		// validate initial output
		initializedOutput := outWire.Signal()
		assert.Equal(t, wire.High, initializedOutput)

		// start proper testing
		highInput := circuit.Change{
			Time:   10,
			Wire:   inWire,
			Signal: wire.High,
		}
		cir.Propagate(highInput)

		newOutput := outWire.Signal()
		assert.Equal(t, wire.Low, newOutput)
	})

	t.Run("AND gate", func(t *testing.T) {
		// set up gate
		cir := circuit.NewCircuit()
		in1 := cir.AddInputWire("input1")
		in2 := cir.AddInputWire("input2")
		outWire := cir.And(in1, in2)

		cir.Initialize()

		// test that High && Low == Low
		high1Input := circuit.Change{
			Time:   10,
			Wire:   in1,
			Signal: wire.High,
		}
		cir.Propagate(high1Input)

		newOutput := outWire.Signal()
		assert.Equal(t, wire.Low, newOutput)

		// test that High && High == High
		high2Input := circuit.Change{
			Time:   10,
			Wire:   in2,
			Signal: wire.High,
		}
		cir.Propagate(high2Input)

		newOutput = outWire.Signal()
		assert.Equal(t, wire.High, newOutput)
	})
}

func TestSingleBitComponents(t *testing.T) {
	t.Run("Half adder", func(t *testing.T) {
		// set up circuit
		cir := circuit.NewCircuit()
		in1 := cir.AddInputWire("input1")
		in2 := cir.AddInputWire("input2")
		sum, carry := cir.HalfAdder(in1, in2)

		testCases := []struct {
			in1Value      wire.Signal
			in2Value      wire.Signal
			expectedSum   wire.Signal
			expectedCarry wire.Signal
		}{
			{
				in1Value: wire.Low,
				in2Value: wire.Low,

				expectedSum:   wire.Low,
				expectedCarry: wire.Low,
			},
			{
				in1Value: wire.Low,
				in2Value: wire.High,

				expectedSum:   wire.High,
				expectedCarry: wire.Low,
			},
			{
				in1Value: wire.High,
				in2Value: wire.Low,

				expectedSum:   wire.High,
				expectedCarry: wire.Low,
			},
			{
				in1Value: wire.High,
				in2Value: wire.High,

				expectedSum:   wire.Low,
				expectedCarry: wire.High,
			},
		}

		for i, tc := range testCases {
			changes := []circuit.Change{
				{
					Time:   i + 1,
					Wire:   in1,
					Signal: tc.in1Value,
				},
				{
					Time:   i + 1,
					Wire:   in2,
					Signal: tc.in2Value,
				},
			}
			for _, ch := range changes {
				cir.Propagate(ch)
			}

			actualSum := sum.Signal()
			actualCarry := carry.Signal()
			assert.EqualValues(t, tc.expectedSum, actualSum)
			assert.EqualValues(t, tc.expectedCarry, actualCarry)
		}
	})
}

// TODO - refactor - add utility functions to break up uint32 into signals?
// TODO - maybe also add utility functions for testing output values more concisely?

func FuzzNot32(f *testing.F) {
	// set up circuit
	cir := circuit.NewCircuit()
	inWire := cir.AddInputWire32("input")
	outWire := cir.Not32(inWire)
	cir.Initialize()

	testcases := []uint32{
		uint32(0),
		^uint32(0),
	}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input uint32) {
		// set input wire values
		changes := []circuit.Change{}

		for i := range 32 {
			ch := circuit.Change{
				Time: 1,
				Wire: inWire.Wire(i),
			}
			if isBitSet(input, i) {
				ch.Signal = wire.High
			} else {
				ch.Signal = wire.Low
			}

			changes = append(changes, ch)
		}

		// push input values into circuit
		cir.Propagate(changes...)

		// read and test output values
		for i := range 32 {
			expectedBitValue := isBitSet(^input, i)
			var expectedSignal wire.Signal
			if expectedBitValue {
				expectedSignal = wire.High
			} else {
				expectedSignal = wire.Low
			}

			actualSignal := outWire.Wire(i).Signal()

			if expectedSignal != actualSignal {
				t.Errorf("Error finding NOT %v, in wire %v: expected %v, actual %v", input, i, expectedSignal, actualSignal)
			}
		}
	})
}

func FuzzAnd32(f *testing.F) {
	// set up circuit
	cir := circuit.NewCircuit()
	inWire1 := cir.AddInputWire32("input1")
	inWire2 := cir.AddInputWire32("input2")
	outWire := cir.And32(inWire1, inWire2)
	cir.Initialize()

	// seed test corpus
	f.Add(uint32(0), uint32(0))
	f.Add(^uint32(0), ^uint32(0))

	f.Fuzz(func(t *testing.T, input1 uint32, input2 uint32) {
		// set input wire values
		changes := []circuit.Change{}
		for i := range 32 {
			ch1 := circuit.Change{
				Time: 1,
				Wire: inWire1.Wire(i),
			}
			if isBitSet(input1, i) {
				ch1.Signal = wire.High
			} else {
				ch1.Signal = wire.Low
			}

			ch2 := circuit.Change{
				Time: 1,
				Wire: inWire2.Wire(i),
			}
			if isBitSet(input2, i) {
				ch2.Signal = wire.High
			} else {
				ch2.Signal = wire.Low
			}

			changes = append(changes, ch1, ch2)
		}

		// push input values into circuit
		cir.Propagate(changes...)

		// read and test output values
		for i := range 32 {
			expectedBitValue := isBitSet(input1&input2, i)
			var expectedSignal wire.Signal
			if expectedBitValue {
				expectedSignal = wire.High
			} else {
				expectedSignal = wire.Low
			}

			actualSignal := outWire.Wire(i).Signal()

			if expectedSignal != actualSignal {
				t.Errorf(
					"Error calculating %v AND %v, in wire %v: expected %v, actual %v",
					input1,
					input2,
					i,
					expectedSignal,
					actualSignal,
				)
			}
		}
	})
}

func FuzzFullAdder32(f *testing.F) {
	// set up circuit
	cir := circuit.NewCircuit()
	inWire1 := cir.AddInputWire32("input1")
	inWire2 := cir.AddInputWire32("input2")
	sumWire, _ := cir.FullAdder32(inWire1, inWire2)
	cir.Initialize()

	// seed test corpus
	f.Add(uint32(0), uint32(0))

	f.Fuzz(func(t *testing.T, input1 uint32, input2 uint32) {
		// set input wire values
		changes := []circuit.Change{}

		for i := range 32 {
			ch1 := circuit.Change{
				Time: 1,
				Wire: inWire1.Wire(i),
			}
			if isBitSet(input1, i) {
				ch1.Signal = wire.High
			} else {
				ch1.Signal = wire.Low
			}

			ch2 := circuit.Change{
				Time: 1,
				Wire: inWire2.Wire(i),
			}
			if isBitSet(input2, i) {
				ch2.Signal = wire.High
			} else {
				ch2.Signal = wire.Low
			}

			changes = append(changes, ch1, ch2)
		}

		// push input values into circuit
		cir.Propagate(changes...)

		// read and test output value
		expectedValue := uint32(input1 + input2)
		actualValue := sumWire.AsUint32()

		if expectedValue != actualValue {
			t.Errorf(
				"Error calculating %v + %v: expected %v, actual %v",
				input1,
				input2,
				expectedValue,
				actualValue,
			)
		}
	})
}

func isBitSet(n uint32, idx int) bool {
	bitmask := uint32(1) << idx
	return (n & bitmask) != 0
}
