package circuit

import "github.com/DylanSp/go-logical-circuits/log"

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
	changes, ok := sim.pendingChanges.GetNextChanges()
	if !ok {
		return Tick{}, false
	}

	// Apply the complete state transition for this timestamp first.
	// TODO - deduplicate identical changes?
	// TODO - how to handle inconsistent changes (the same wire getting set to different values at the same tick)?
	// maybe just panic if there are inconsistent changes, since that indicates a non-well-formed circuit?
	for _, change := range changes {
		log.Logf("Executing change at t=%v, changing wire %v to %v\n", change.Time, change.Wire, change.Signal)
		change.Wire.SetSignal(change.Signal)
	}

	// Calculate downstream effects only after all simultaneous changes apply.
	for _, change := range changes {
		for _, wire := range sim.circuit.wires {
			sim.Schedule(wire.HandleChange(change)...)
		}
	}

	return Tick{
		Time:    changes[0].Time,
		Changes: changes,
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
	for inputWire := range sim.circuit.inputWires {
		sim.Schedule(Change{
			Time:   0,
			Wire:   inputWire,
			Signal: Low,
		})
	}

	sim.RunUntilStable()
}
