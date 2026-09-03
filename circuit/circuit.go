package circuit

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit/wire"
)

// single change in a single wire
type Change struct {
	Time int

	Wire   *wire.Wire
	Signal wire.Signal
}

func (ch Change) String() string {
	return fmt.Sprintf("%v (t=%v) new value: %v", ch.Wire, ch.Time, ch.Signal)
}

// processing delay of a single gate (of any time)
const gateDelay = 2

type Component struct {
	// used for checking if a change affects this component
	// uses a set to enforce uniqueness
	inputWires map[*wire.Wire]struct{}

	// internal logic for responding to changes
	// does *not* update wires' states; Circuit.ExecuteChanges() will take the changes from this and update states appropriately
	processChangeLogic func(Change) []Change
}

func (com Component) HandleChange(ch Change) []Change {
	if _, ok := com.inputWires[ch.Wire]; !ok {
		// the change isn't one of this component's input wires; no changes produced
		return []Change{}
	}

	return com.processChangeLogic(ch)
}

type Circuit struct {
	wires      map[string]*wire.Wire // wires by name
	components []Component

	internalWireCount int // used for giving internal wires unique names
}

func NewCircuit() Circuit {
	return Circuit{
		wires:             map[string]*wire.Wire{},
		components:        []Component{},
		internalWireCount: 0,
	}
}

func (cir *Circuit) addComponents(components ...Component) {
	cir.components = append(cir.components, components...)
}

// return the added wire so callers can refer to it for setting up components
func (cir *Circuit) AddWire(wireName string) *wire.Wire {
	_, hasWire := cir.wires[wireName]
	if hasWire {
		panic(fmt.Sprintf("Wire %v already present", wireName))
	}

	newWire := wire.New(wireName)
	cir.wires[wireName] = newWire

	return newWire
}

// used for adding internal wires inside compound components
func (cir *Circuit) addInternalWire() *wire.Wire {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	wire := cir.AddWire(wireName)
	cir.internalWireCount++
	return wire
}

/*****
Overall flow of changes:
1. a Change to an input wire is created and passed to Propagate()
2. Propagate() adds this as the initial change to agenda
3. Propagate() pops the next Change(s) to be made from the agenda, passes it/them to executeChange()
4. executeChange() updates the change's wire with the new signal from the change
5. executeChange() calls .HandleChange() on all components to see if any new changes are created
6. Any component with an input wire whose signal was changed creates a new Change with the updated value for output wires, returns it from HandleChange()
7. All Changes created this way are collected by executeChange() and passed back to Propagate(), which adds them to the agenda
8. Propagate() returns to step 3; this loops until there are no more downstream changes to execute

In the simple case of a single gate:
1. Change is created affecting an input wire
2. First iteration of executeChange() updates that wire's value
3. First iteration of executeChange() calls .HandleChange(); gate calculates the new value for its output wire and returns a Change specifying that
4. Second iteration of executeChange() receives the Change for the output wire and updates its value
*/

// execute a single change, possibly producing further changes
// the only wire whose signal is actually changed is ch.Wire
func (cir *Circuit) executeChange(ch Change) []Change {
	fmt.Printf("Executing change at t=%v, Changing wire %v to %v\n", ch.Time, ch.Wire, ch.Signal)

	if _, ok := cir.wires[ch.Wire.Name()]; !ok {
		panic(fmt.Sprintf("Circuit doesn't contain wire %v", ch.Wire)) // TODO - should we have/do we need this check?
	}

	// update wire from change
	// fmt.Printf("Calling SetSignal on wire %v with value %v\n", ch.Wire, ch.Signal)
	ch.Wire.SetSignal(wire.Signal(ch.Signal))

	// now check for downstream changes
	downstreamChanges := []Change{}
	for _, com := range cir.components {
		changeResults := com.HandleChange(ch)
		for _, chResult := range changeResults {
			downstreamChanges = append(downstreamChanges, chResult)
		}
	}

	return downstreamChanges
}

// propagate a single change through the circuit until it stabilizes
// TODO - how to handle cases where circuit never stabilizes?
// TODO - return an iterator of some sort?
func (cir *Circuit) Propagate(initialChange Change) {
	agenda := NewChangeQueue()
	agenda.AddChange(initialChange)

	fmt.Printf("Propagating from change at t=%v\n", initialChange.Time)

	for agenda.Length() > 0 {
		nextChanges, _ := agenda.GetNextChanges()

		// TODO - how to handle the case where there is >1 nextChange at the same time?
		// TODO - check for possibility of overlap and resolve/combine somehow?
		for _, nextChange := range nextChanges {
			// TODO - having this check messes with circuit initialization
			// if prevSignal == wire.Signal(newSignal) {
			// 	continue // no actual change in signal => no-op
			// }

			downstreamChanges := cir.executeChange(nextChange)
			for _, downstream := range downstreamChanges {
				agenda.AddChange(downstream)
			}
		}
	}

	fmt.Printf("Done propagating\n\n")
}
