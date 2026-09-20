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
			bitmask := uint32(1) << i
			ithBitIsSet := (input & bitmask) != 0

			ch := circuit.Change{
				Time: 1,
				Wire: inWire.Wire(i),
			}
			if ithBitIsSet {
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
			bitmask := uint32(1) << i
			expectedBitValue := (input & bitmask) == 0 // bits set in input should be unset in output
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
