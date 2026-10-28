package circuit_test

import (
	"testing"

	"github.com/DylanSp/go-logical-circuits/circuit"
	"github.com/stretchr/testify/assert"
)

func TestTick(t *testing.T) {
	t.Run("Tick only steps forward one time step, even if next scheduled change is multiple time steps ahead", func(t *testing.T) {
		// set up simple circuit
		cir := circuit.NewCircuit()
		sim := circuit.NewSimulation(cir)
		inWire := cir.AddInputWire("input")
		_ = cir.Noop(inWire)

		sim.Initialize()

		sim.ScheduleAhead(inWire, circuit.High, 50)

		// even though there's a change scheduled, this tick shouldn't have any changes
		tick, changesHappened := sim.Tick()

		assert.False(t, changesHappened)
		assert.Len(t, tick.Changes, 0)
	})
}
