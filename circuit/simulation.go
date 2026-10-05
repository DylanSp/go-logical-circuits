package circuit

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/DylanSp/go-logical-circuits/log"
	"github.com/DylanSp/go-logical-circuits/metrics"
)

// Tick describes all changes applied during one simulation tick.
type Tick struct {
	Time    int
	Changes []Change
}

// Simulation contains the runtime state for executing a Circuit.
type Simulation struct {
	circuit        *Circuit
	pendingChanges *ChangeQueue
}

func NewSimulation(circuit *Circuit) *Simulation {
	return &Simulation{
		circuit:        circuit,
		pendingChanges: NewChangeQueue(),
	}
}

// Schedule adds changes to the simulation's pending agenda.
func (sim *Simulation) Schedule(changes ...Change) {
	for _, change := range changes {
		sim.pendingChanges.AddChange(change)
	}
}

// IsStable reports whether the simulation has no pending changes.
func (sim *Simulation) IsStable() bool {
	return sim.pendingChanges.Length() == 0
}

/*****
Overall flow of changes:
1. a Change to an input wire is created and added to the agenda with Schedule()
2. Tick() pops the next Change(s) to be made from the agenda, applies them to the relevant wires
3. Tick() then calls .HandleChange() on all wires to see if any new changes are created
4. Any wires with an input wire whose signal was changed create a new Change with the updated value for output wires, returns it from HandleChange()
5. All Changes created this way are scheduled on the agenda from within Tick()
6. Tick() then returns all changes that happened during the tick

RunUntilStable() calls Tick() in a loop until there are no more changes pending.

In the simple case of a single gate:
1. A Change is created affecting an input wire and scheduled
2. First call to Tick() updates that wire's value
3. First call to Tick() also calls wire.HandleChange();
	gate calculates the new value for its output wire and returns a Change specifying that new value
4. The Change returned from wire.HandleChange() is scheduled on the agenda
5. Second call to Tick() gets the Change for the output wire from the agenda
6. Second call to Tick() updates the output wire's value
*/

// Tick advances the simulation to process the next batch of changes.
// All simultaneous changes are applied before their downstream effects are
// calculated and scheduled.
func (sim *Simulation) Tick() (Tick, bool) {
	return sim.tick(false)
}

func (sim *Simulation) tick(isInitialization bool) (Tick, bool) {
	startTime := time.Now()

	changes, ok := sim.pendingChanges.GetNextChanges()
	if !ok {
		// no pending changes this tick, no metrics to record

		return Tick{}, false
	}

	// only apply changes that actually change the state of a wire
	var actualChanges []Change
	if true /* isInitialization */ {
		// on initialization, don't ignore *any* changes;
		// all wires start Low, we need to drive initial changes through to make circuit consistent
		actualChanges = changes
	} else {
		for _, ch := range changes {
			if ch.Signal != ch.Wire.Signal() {
				actualChanges = append(actualChanges, ch)
			}
		}
	}

	if len(actualChanges) == 0 {
		duration := time.Since(startTime)
		metrics.RecordTick(changes[0].Time, duration, 0)

		return Tick{}, false
	}

	// check for inconsistent changes
	changesByWireName := map[string][]Change{}
	for _, ch := range actualChanges {
		changesByWireName[ch.Wire.Name()] = append(changesByWireName[ch.Wire.Name()], ch)
	}

	for wire, wireChanges := range changesByWireName {
		if len(wireChanges) < 2 {
			continue
		}

		initialChangeValue := wireChanges[0].Signal
		for _, wireChange := range wireChanges[1:] {
			if wireChange.Signal != initialChangeValue {
				panic(fmt.Sprintf("Conflicting changes for wire %v at time %v; aborting", wire, wireChange.Time))
			}
		}
	}

	// Apply the complete state transition for this timestamp first.
	// TODO - deduplicate identical changes?
	// TODO - how to handle inconsistent changes (the same wire getting set to different values at the same tick)?
	// maybe just panic if there are inconsistent changes, since that indicates a non-well-formed circuit?
	for _, ch := range actualChanges {
		log.Logf("Executing change at t=%v, changing wire %v to %v\n", ch.Time, ch.Wire, ch.Signal)
		ch.Wire.SetSignal(ch.Signal)
	}

	// Calculate downstream effects only after all simultaneous changes apply.
	for _, ch := range actualChanges {
		for _, wire := range sim.circuit.wires {
			downstreamChanges := wire.HandleChange(ch)

			// optimization - if a downstream change wouldn't actually change the value of a wire,
			// (because that wire is already set to that value)
			// ignore it
			// TODO - check for conflicting changes here?
			// if we don't check, conflicting changes coming from the same source will be ignored;
			// the changes matching the wire's current value will be filtered out,
			// while the changes that don't match will be scheduled
			// However, this shouldn't be very likely,
			// a single component would have to have conflicting changes from one call to processChangeLogic()
			for _, downstream := range downstreamChanges {
				if downstream.Signal != downstream.Wire.Signal() {
					sim.Schedule(downstream)
				}
			}

			// sim.Schedule(wire.HandleChange(ch)...)
		}
	}

	duration := time.Since(startTime)
	metrics.RecordTick(actualChanges[0].Time, duration, len(actualChanges))

	return Tick{
		Time:    actualChanges[0].Time,
		Changes: actualChanges,
	}, true
}

// RunUntilStable processes ticks until no changes remain. It will never return
// for circuits which continually schedule future changes, such as clocks
func (sim *Simulation) RunUntilStable() {
	for {
		_, ok := sim.Tick()
		if !ok {
			return
		}
	}
}

// Initialize drives every circuit input Low and propagates those changes.
func (sim *Simulation) Initialize() {
	var initialWires []*Wire

	if true /* len(sim.circuit.inputWires) != 0 */ {
		// if we have input wires set, use them
		initialWires = slices.Collect(maps.Keys(sim.circuit.inputWires))
	} else {
		// no input wires are set

		// if we don't have any wires at all, bail out
		if len(sim.circuit.wires) == 0 {
			return
		}

		// pick a random wire and start with that
		allWires := slices.Collect(maps.Values(sim.circuit.wires))
		initialWires = []*Wire{
			allWires[rand.IntN(len(allWires))],
		}
	}

	for _, inputWire := range initialWires {
		sim.Schedule(Change{
			Time:   0,
			Wire:   inputWire,
			Signal: Low,
		})
	}

	for {
		_, ok := sim.tick(true)
		if !ok {
			return
		}
	}
}
