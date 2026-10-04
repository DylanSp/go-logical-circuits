package circuit_test

import (
	"testing"

	"github.com/DylanSp/go-logical-circuits/circuit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimulation(t *testing.T) {
	t.Run("Tick applies simultaneous changes before calculating effects", func(t *testing.T) {
		cir := circuit.NewCircuit()
		in1 := cir.AddInputWire("input1")
		in2 := cir.AddInputWire("input2")
		out := cir.And(in1, in2)
		sim := circuit.NewSimulation(cir)
		sim.Initialize()

		sim.Schedule(
			circuit.Change{Time: 10, Wire: in1, Signal: circuit.High},
			circuit.Change{Time: 10, Wire: in2, Signal: circuit.High},
		)

		inputTick, ok := sim.Tick()

		require.True(t, ok)
		assert.Equal(t, 10, inputTick.Time)

		// check that Tick() returned *both* changes processed
		assert.Len(t, inputTick.Changes, 2)

		// check that downstream changes have *not* yet been processed;
		// once they have, `out` will be High
		assert.Equal(t, circuit.Low, out.Signal())

		_, ok = sim.Tick()

		require.True(t, ok)

		// check that downsteam changes have now been processed and `out` is driven to High
		assert.Equal(t, circuit.High, out.Signal())
		assert.True(t, sim.IsStable())

		_, ok = sim.Tick()
		assert.False(t, ok)
	})
}
