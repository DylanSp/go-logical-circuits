package circuit

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit/wire"
)

// single change in a single wire
type Change struct {
	Time   int
	Wire   *wire.Wire
	Signal bool // TODO - change to proper Signal type?
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

// TODO - put functions for making gates in another file?

func MkNotGate(inWire *wire.Wire, outWire *wire.Wire) Component {
	gate := Component{
		inputWires: map[*wire.Wire]struct{}{},
	}
	gate.inputWires[inWire] = struct{}{}

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   outWire,
				Signal: !ch.Signal,
			},
		}
	}

	return gate
}

func MkAndGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) Component {
	gate := Component{
		inputWires: map[*wire.Wire]struct{}{},
	}
	gate.inputWires[in1] = struct{}{}
	gate.inputWires[in2] = struct{}{}

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time: ch.Time + gateDelay,
				Wire: out,

				// TODO - incorrect - this looks up old signals instead of taking into account new values from the change
				Signal: bool(in1.Signal()) && bool(in2.Signal()),
			},
		}
	}

	return gate
}

type Circuit struct {
	wires      map[string]*wire.Wire // wires by name
	components []Component

	internalWireCount int
}

func NewCircuit() Circuit {
	return Circuit{
		wires:             map[string]*wire.Wire{},
		components:        []Component{},
		internalWireCount: 0,
	}
}

func (cir *Circuit) AddComponents(components ...Component) {
	// TODO - do we need to do anything else?
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

// used for adding internal wires inside components
// TODO - does this need to be exported?
func (cir *Circuit) AddInternalWire() {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	cir.AddWire(wireName)
	cir.internalWireCount++
}

// execute a single change and its immediate effects, possibly producing further changes
func (cir *Circuit) executeChange(ch Change) []Change {
	fmt.Printf("Executing change at t=%v, Changing wire %v to %v\n", ch.Time, ch.Wire, ch.Signal)

	if _, ok := cir.wires[ch.Wire.Name()]; !ok {
		panic(fmt.Sprintf("Circuit doesn't contain wire %v", ch.Wire)) // TODO - should we have this check?
	}

	downstreamChanges := []Change{}
	for _, com := range cir.components {
		changeResults := com.HandleChange(ch)
		for _, chResult := range changeResults {
			fmt.Printf("Calling SetSignal on wire %v with value %v\n", chResult.Wire, chResult.Signal)
			chResult.Wire.SetSignal(wire.Signal(chResult.Signal))

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
