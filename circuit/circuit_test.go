package circuit_test

import (
	"fmt"
	"testing"

	"github.com/DylanSp/go-logical-circuits/bits"
	"github.com/DylanSp/go-logical-circuits/circuit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		sim.Schedule(inWire, circuit.High)

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
		sim.Schedule(in1, circuit.High)

		sim.RunUntilStable()

		newOutput := outWire.Signal()
		assert.Equal(t, circuit.Low, newOutput)

		// test that High && High == High
		sim.Schedule(in2, circuit.High)
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

		for _, tc := range testCases {
			sim.Schedule(in1, tc.in1Value)
			sim.Schedule(in2, tc.in2Value)

			sim.RunUntilStable()

			actualSum := sum.Signal()
			actualCarry := carry.Signal()
			assert.EqualValues(t, tc.expectedSum, actualSum)
			assert.EqualValues(t, tc.expectedCarry, actualCarry)
		}
	})

	t.Run("Flip-flop", func(t *testing.T) {
		// set up circuit
		cir := circuit.NewCircuit()
		input := cir.AddInputWire("input")

		clock := cir.AddInputWire("clock") // will be manually controlled for testing instead of an actual clock
		output := cir.FlipFlop(input, clock)

		sim := circuit.NewSimulation(cir)
		sim.Initialize()

		// ensure input is set to Low,
		// run clock through a full cycle to stabilize flip-flop in a known & valid state (storing & outputting Low)
		sim.Schedule(input, circuit.Low)
		sim.ScheduleAhead(clock, circuit.Low, 5)
		sim.ScheduleAhead(clock, circuit.High, 10)
		sim.ScheduleAhead(clock, circuit.Low, 15)

		// let the initial values propagate until circuit stabilizes
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.Low, output.Signal(), "Output should be Low after stabilizing clock and running 100 ticks")

		// set input to high; should be ignored by flip-flop, because clock remains low
		sim.Schedule(input, circuit.High)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.Low, output.Signal(), "Output should still be Low after setting input High while clock remains Low")

		// set clock to high (clock raising edge); flip-flop should now accept and store High input
		sim.Schedule(clock, circuit.High)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should be High after clock raising edge with High input")

		// set input to low;
		// should be ignored by flip-flop and output should remain high, because clock hasn't changed (still high)
		sim.Schedule(input, circuit.Low)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should still be High after setting input Low while clock remains High")

		// set clock to low (clock falling edge); output should remain high
		sim.Schedule(clock, circuit.Low)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should still be High on clock falling edge, even with input Low")

		// set input to low; should be ignored by flip-flop (remaining high), because clock is still low
		sim.Schedule(input, circuit.Low)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.High, output.Signal(), "Output should still be High after setting input Low while clock remains Low")

		sim.Schedule(clock, circuit.High)
		for range 100 {
			sim.Tick()
		}

		assert.EqualValues(t, circuit.Low, output.Signal(), "Output should be Low after clock raising edge with Low input")
	})
}

func Test32BitComponents(t *testing.T) {
	t.Run("Single-bit multiplexer", func(t *testing.T) {
		cir := circuit.NewCircuit()
		inWire0 := cir.AddInputWire32("in0")
		inWire1 := cir.AddInputWire32("in1")
		selector := cir.AddInputWire("sel")
		mux := cir.Mux32(inWire1, inWire0, selector)

		sim := circuit.NewSimulation(cir)
		sim.Initialize()

		input0Value := uint32(0)
		input1Value := uint32(1)

		sim.Schedule32(inWire0, input0Value)
		sim.Schedule32(inWire1, input1Value)

		// set selector to 0; output should be input0
		sim.Schedule(selector, circuit.Low)

		sim.RunUntilStable()

		assert.EqualValues(t, input0Value, mux.AsUint32())

		// set selector to 1; output should be input1
		sim.Schedule(selector, circuit.High)

		sim.RunUntilStable()

		assert.EqualValues(t, input1Value, mux.AsUint32())
	})

	t.Run("N-bit multiplexer", func(t *testing.T) {
		cir := circuit.NewCircuit()
		inWire0 := cir.AddInputWire32("in0")
		inWire1 := cir.AddInputWire32("in1")
		inWire2 := cir.AddInputWire32("in2")
		inWire3 := cir.AddInputWire32("in3")
		inputWires := []*circuit.Wire32{
			inWire0,
			inWire1,
			inWire2,
			inWire3,
		}

		selector0 := cir.AddInputWire("sel0")
		selector1 := cir.AddInputWire("sel1") // most-significant bit of selector

		// Mux32N takes selectors as a big-endian slice
		selectors := []*circuit.Wire{
			selector1,
			selector0,
		}

		mux, ok := cir.Mux32N(inputWires, selectors)
		require.True(t, ok)

		sim := circuit.NewSimulation(cir)
		sim.Initialize()

		inputValues := []uint32{
			10,
			1,
			2,
			3,
		}

		testCases := []struct {
			selector1Value  circuit.Signal
			selector0Value  circuit.Signal
			selectedWireNum int
		}{
			{
				selector1Value:  circuit.Low,
				selector0Value:  circuit.Low,
				selectedWireNum: 0,
			},
			{
				selector1Value:  circuit.Low,
				selector0Value:  circuit.High,
				selectedWireNum: 1,
			},
			{
				selector1Value:  circuit.High,
				selector0Value:  circuit.Low,
				selectedWireNum: 2,
			},
			{
				selector1Value:  circuit.High,
				selector0Value:  circuit.High,
				selectedWireNum: 3,
			},
		}

		for _, tc := range testCases {
			tcName := fmt.Sprintf("selector=%v", tc.selectedWireNum)
			t.Run(tcName, func(t *testing.T) {
				// inputChanges := []circuit.Change{}

				for inputNum := range 4 {
					// singleInputChanges := inputWires[inputNum].ChangesFromUint32(inputValues[inputNum], i+1)
					// inputChanges = append(inputChanges, singleInputChanges...)

					sim.Schedule32(inputWires[inputNum], inputValues[inputNum])
				}

				sim.Schedule(selector0, tc.selector0Value)
				sim.Schedule(selector1, tc.selector1Value)

				sim.RunUntilStable()

				outputValue := mux.AsUint32()
				assert.EqualValues(t, inputWires[tc.selectedWireNum].AsUint32(), outputValue)
			})
		}
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

	sim.Schedule(in2, circuit.High)
	sim.RunUntilStable()

	// remove in1, add in3 from AND gate's inputs
	gate.SetInputs(in2, in3)

	// if SetInputs correctly reset gate's inputs, this should drive gate High
	sim.Schedule(in3, circuit.High)
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
		sim.Schedule32(inWire, input)
		sim.RunUntilStable()

		// read and test output values
		for i := range 32 {
			expectedSignal := circuit.Signal(bits.IsBitSet(^input, i))
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
		sim.Schedule32(inWire1, input1)
		sim.Schedule32(inWire2, input2)
		sim.RunUntilStable()

		// read and test output values
		for i := range 32 {
			expectedSignal := circuit.Signal(bits.IsBitSet(input1&input2, i))
			actualSignal := outWire.Wire(i).Signal()

			if expectedSignal != actualSignal {
				t.Errorf(
					"Error calculating %v AND %v, in wire %v: expected %v, actual %v",
					inWire1,
					inWire2,
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
		sim.Schedule32(inWire1, input1)
		sim.Schedule32(inWire2, input2)
		sim.RunUntilStable()

		// read and test output value
		expectedValue := uint32(input1 + input2)
		actualValue := sumWire.AsUint32()

		if expectedValue != actualValue {
			t.Errorf(
				"Error calculating %v + %v: expected %v, actual %v",
				inWire1,
				inWire2,
				expectedValue,
				actualValue,
			)
		}

		expectedOverflow := (uint64(input1) + uint64(input2)) != uint64(expectedValue)
		actualOverflow := overflow.IsHigh()

		if expectedOverflow != actualOverflow {
			t.Errorf(
				"Error detecting overflow when calculating %v + %v: expected %v, actual %v",
				inWire1,
				inWire2,
				expectedOverflow,
				actualOverflow,
			)
		}
	})
}
