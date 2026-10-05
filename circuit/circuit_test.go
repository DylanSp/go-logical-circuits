package circuit_test

import (
	"testing"

	"github.com/DylanSp/go-logical-circuits/circuit"
	"github.com/stretchr/testify/assert"
)

func TestBasicGates(t *testing.T) {
	t.Run("NOT gate", func(t *testing.T) {
		// set up gate
		cir := circuit.NewCircuit()
		sim := circuit.NewSimulation(cir)
		inWire := cir.AddInputWire("input")
		outWire := cir.Not(inWire)

		sim.Initialize()

		// validate initial output
		initializedOutput := outWire.Signal()
		assert.Equal(t, circuit.High, initializedOutput)

		// start proper testing
		highInput := circuit.Change{
			Time:   10,
			Wire:   inWire,
			Signal: circuit.High,
		}
		sim.Schedule(highInput)
		sim.RunUntilStable()

		newOutput := outWire.Signal()
		assert.Equal(t, circuit.Low, newOutput)
	})

	t.Run("AND gate", func(t *testing.T) {
		// set up gate
		cir := circuit.NewCircuit()
		sim := circuit.NewSimulation(cir)
		in1 := cir.AddInputWire("input1")
		in2 := cir.AddInputWire("input2")
		outWire := cir.And(in1, in2)

		sim.Initialize()

		// test that High && Low == Low
		high1Input := circuit.Change{
			Time:   10,
			Wire:   in1,
			Signal: circuit.High,
		}
		sim.Schedule(high1Input)
		sim.RunUntilStable()

		newOutput := outWire.Signal()
		assert.Equal(t, circuit.Low, newOutput)

		// test that High && High == High
		high2Input := circuit.Change{
			Time:   10,
			Wire:   in2,
			Signal: circuit.High,
		}
		sim.Schedule(high2Input)
		sim.RunUntilStable()

		newOutput = outWire.Signal()
		assert.Equal(t, circuit.High, newOutput)
	})
}

func TestSingleBitComponents(t *testing.T) {
	t.Run("Half adder", func(t *testing.T) {
		// set up circuit
		cir := circuit.NewCircuit()
		sim := circuit.NewSimulation(cir)
		in1 := cir.AddInputWire("input1")
		in2 := cir.AddInputWire("input2")
		sum, carry := cir.HalfAdder(in1, in2)

		testCases := []struct {
			in1Value      circuit.Signal
			in2Value      circuit.Signal
			expectedSum   circuit.Signal
			expectedCarry circuit.Signal
		}{
			{
				in1Value: circuit.Low,
				in2Value: circuit.Low,

				expectedSum:   circuit.Low,
				expectedCarry: circuit.Low,
			},
			{
				in1Value: circuit.Low,
				in2Value: circuit.High,

				expectedSum:   circuit.High,
				expectedCarry: circuit.Low,
			},
			{
				in1Value: circuit.High,
				in2Value: circuit.Low,

				expectedSum:   circuit.High,
				expectedCarry: circuit.Low,
			},
			{
				in1Value: circuit.High,
				in2Value: circuit.High,

				expectedSum:   circuit.Low,
				expectedCarry: circuit.High,
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
				sim.Schedule(ch)
				sim.RunUntilStable()
			}

			actualSum := sum.Signal()
			actualCarry := carry.Signal()
			assert.EqualValues(t, tc.expectedSum, actualSum)
			assert.EqualValues(t, tc.expectedCarry, actualCarry)
		}
	})

	t.Run("Flip-flop", func(t *testing.T) {
		// IMPORTANT NOTE FOR DEBUG/TEST CODE:
		// if manually using .Tick(), make sure the number of Tick() calls matches up with the time on scheduled changes
		// need to run for 100 Ticks before changes scheduled at Time=100 will be executed!

		// set up circuit
		cir := circuit.NewCircuit()
		input := cir.AddInputWire("input")

		clock := cir.AddInputWire("clock") // will be manually controlled for testing instead of an actual clock
		output := cir.FlipFlop(input, clock)

		sim := circuit.NewSimulation(cir)
		sim.Initialize()

		// ensure input is set to Low,
		// run clock through a full cycle to stabilize flip-flop in a known & valid state (storing & outputting Low)
		sim.Schedule(
			circuit.Change{
				Time:   0,
				Wire:   input,
				Signal: circuit.Low,
			},

			circuit.Change{
				Time:   5,
				Wire:   clock,
				Signal: circuit.Low,
			},
			circuit.Change{
				Time:   10,
				Wire:   clock,
				Signal: circuit.High,
			},
			circuit.Change{
				Time:   15,
				Wire:   clock,
				Signal: circuit.Low,
			},
		)

		// let the initial values propagate until circuit stabilizes
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.Low, output.Signal(), "Output should be Low after stabilizing clock and running 100 ticks")

		// set input to high; should be ignored by flip-flop, because clock remains low
		sim.Schedule(
			circuit.Change{
				Time:   100,
				Wire:   input,
				Signal: circuit.High,
			},
		)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.Low, output.Signal(), "Output should still be Low after setting input High while clock remains Low")

		// set clock to high (clock raising edge); flip-flop should now accept and store High input
		sim.Schedule(
			circuit.Change{
				Time:   200,
				Wire:   clock,
				Signal: circuit.High,
			},
		)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should be High after clock raising edge with High input")

		// set input to low;
		// should be ignored by flip-flop and output should remain high, because clock hasn't changed (still high)
		sim.Schedule(
			circuit.Change{
				Time:   300,
				Wire:   input,
				Signal: circuit.Low,
			},
		)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should still be High after setting input Low while clock remains High")

		// set clock to low (clock falling edge); output should remain high
		sim.Schedule(
			circuit.Change{
				Time:   400,
				Wire:   clock,
				Signal: circuit.Low,
			})

		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should still be High on clock falling edge, even with input Low")

		// set input to low; should be ignored by flip-flop (remaining high), because clock is still low
		sim.Schedule(
			circuit.Change{
				Time:   500,
				Wire:   input,
				Signal: circuit.Low,
			},
		)

		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should still be High after setting input Low while clock remains Low")

		sim.Schedule(
			circuit.Change{
				Time:   600,
				Wire:   clock,
				Signal: circuit.High,
			},
		)

		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.Low, output.Signal(), "Output should be Low after clock raising edge with Low input")
	})
}

func TestSetInputs(t *testing.T) {
	cir := circuit.NewCircuit()
	in1 := cir.AddInputWire("input1")
	in2 := cir.AddInputWire("input2")
	in3 := cir.AddInputWire("input3")
	gate := cir.And(in1, in2)

	sim := circuit.NewSimulation(cir)
	sim.Initialize()
	sim.RunUntilStable()

	high2Input := circuit.Change{
		Time:   10,
		Wire:   in2,
		Signal: circuit.High,
	}
	sim.Schedule(high2Input)
	sim.RunUntilStable()

	// remove in1, add in3 from AND gate's inputs
	gate.SetInputs(in2, in3)

	// if SetInputs correctly reset gate's inputs, this should drive gate High
	high3Input := circuit.Change{
		Time:   20,
		Wire:   in3,
		Signal: circuit.High,
	}
	sim.Schedule(high3Input)
	sim.RunUntilStable()

	assert.Equal(t, circuit.High, gate.Signal())
}

func FuzzNot32(f *testing.F) {
	// set up circuit
	cir := circuit.NewCircuit()
	sim := circuit.NewSimulation(cir)
	inWire := cir.AddInputWire32("input")
	outWire := cir.Not32(inWire)
	sim.Initialize()

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
				Time:   1,
				Wire:   inWire.Wire(i),
				Signal: circuit.Signal(isBitSet(input, i)),
			}

			changes = append(changes, ch)
		}

		// push input values into circuit
		sim.Schedule(changes...)
		sim.RunUntilStable()

		// read and test output values
		for i := range 32 {
			expectedSignal := circuit.Signal(isBitSet(^input, i))
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
	sim := circuit.NewSimulation(cir)
	inWire1 := cir.AddInputWire32("input1")
	inWire2 := cir.AddInputWire32("input2")
	outWire := cir.And32(inWire1, inWire2)
	sim.Initialize()

	// seed test corpus
	f.Add(uint32(0), uint32(0))
	f.Add(^uint32(0), ^uint32(0))

	f.Fuzz(func(t *testing.T, input1 uint32, input2 uint32) {
		// set input wire values
		changes := []circuit.Change{}
		for i := range 32 {
			ch1 := circuit.Change{
				Time:   1,
				Wire:   inWire1.Wire(i),
				Signal: circuit.Signal(isBitSet(input1, i)),
			}
			ch2 := circuit.Change{
				Time:   1,
				Wire:   inWire2.Wire(i),
				Signal: circuit.Signal(isBitSet(input2, i)),
			}
			changes = append(changes, ch1, ch2)
		}

		// push input values into circuit
		sim.Schedule(changes...)
		sim.RunUntilStable()

		// read and test output values
		for i := range 32 {
			expectedSignal := circuit.Signal(isBitSet(input1&input2, i))
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
	sim := circuit.NewSimulation(cir)
	inWire1 := cir.AddInputWire32("input1")
	inWire2 := cir.AddInputWire32("input2")
	sumWire, overflow := cir.FullAdder32(inWire1, inWire2)
	sim.Initialize()

	// seed test corpus
	f.Add(uint32(0), uint32(0))
	f.Add(^uint32(0), uint32(1))

	f.Fuzz(func(t *testing.T, input1 uint32, input2 uint32) {
		// set input wire values
		changes := []circuit.Change{}

		for i := range 32 {
			ch1 := circuit.Change{
				Time:   1,
				Wire:   inWire1.Wire(i),
				Signal: circuit.Signal(isBitSet(input1, i)),
			}
			ch2 := circuit.Change{
				Time:   1,
				Wire:   inWire2.Wire(i),
				Signal: circuit.Signal(isBitSet(input2, i)),
			}
			changes = append(changes, ch1, ch2)
		}

		// push input values into circuit
		sim.Schedule(changes...)
		sim.RunUntilStable()

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

		expectedOverflow := (uint64(input1) + uint64(input2)) != uint64(expectedValue)
		actualOverflow := overflow.IsHigh()

		if expectedOverflow != actualOverflow {
			t.Errorf(
				"Error detecting overflow when calculating %v + %v: expected %v, actual %v",
				input1,
				input2,
				expectedOverflow,
				actualOverflow,
			)
		}
	})
}

func isBitSet(n uint32, idx int) bool {
	bitmask := uint32(1) << idx
	return (n & bitmask) != 0
}
