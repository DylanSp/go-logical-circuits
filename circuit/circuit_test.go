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
		outWire := cir.AddOutputWire("output")
		cir.MkNotGate(inWire, outWire)

		// initialize input wire
		initialization := circuit.Change{
			Time:   0,
			Wire:   inWire,
			Signal: wire.Low,
		}
		cir.Propagate(initialization)

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
		outWire := cir.AddOutputWire("output")
		cir.MkAndGate(in1, in2, outWire)

		// initialize input wires
		initializations := []circuit.Change{
			{
				Time:   0,
				Wire:   in1,
				Signal: wire.Low,
			},
			{
				Time:   1,
				Wire:   in2,
				Signal: wire.Low,
			},
		}
		for _, init := range initializations {
			cir.Propagate(init)
		}

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
