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
	// TODO - instead of tracking inputWires separately, initialize by setting a single wire to Low, then propagating?
	// TODO - would allow just having a single AddWire() method instead of Add{Input,Output,Auxiliary}Wire/
	inputWires map[*wire.Wire]struct{} // tracked for initialization
	wires      map[string]*wire.Wire   // all wires by name; tracked to avoid duplicate wires

	components []Component

	internalWireCount int // used for giving internal wires unique names
}

func NewCircuit() Circuit {
	return Circuit{
		wires:             map[string]*wire.Wire{},
		inputWires:        map[*wire.Wire]struct{}{},
		components:        []Component{},
		internalWireCount: 0,
	}
}

// return the added wire so callers can refer to it for setting up components
func (cir *Circuit) addWire(wireName string) *wire.Wire {
	_, hasWire := cir.wires[wireName]
	if hasWire {
		panic(fmt.Sprintf("Wire %v already present", wireName))
	}

	newWire := wire.New(wireName)
	cir.wires[wireName] = newWire

	return newWire
}

func (cir *Circuit) addWire32(wireName string) *wire.Wire32 {
	// TODO - check if wire is already present? (circuit would need another field to track Wire32's)

	newWires := [32]*wire.Wire{}
	for i := range 32 {
		newWires[i] = cir.addWire(fmt.Sprintf("%v-%v", wireName, i))
	}
	return wire.FromWires(wireName, newWires)
}

func (cir *Circuit) AddInputWire(wireName string) *wire.Wire {
	wire := cir.addWire(wireName)
	cir.inputWires[wire] = struct{}{}
	return wire
}

func (cir *Circuit) AddInputWire32(wireName string) *wire.Wire32 {
	wire := cir.addWire32(wireName)

	// add the single-bit wires as input wires to the circuit
	for i := range 32 {
		cir.inputWires[wire.Wire(i)] = struct{}{}
	}

	return wire
}

func (cir *Circuit) AddOutputWire(wireName string) *wire.Wire {
	wire := cir.addWire(wireName)
	return wire
}

// auxiliary wires - neither inputs nor outputs
func (cir *Circuit) AddAuxiliaryWire(wireName string) *wire.Wire {
	wire := cir.addWire(wireName)
	return wire
}

// used for adding internal wires inside compound components
func (cir *Circuit) addInternalWire() *wire.Wire {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	wire := cir.addWire(wireName)
	cir.internalWireCount++
	return wire
}

// TODO - do we need this?
func (cir *Circuit) addInternalWire32() *wire.Wire32 {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	wire := cir.addWire32(wireName)
	cir.internalWireCount++
	return wire
}

func (cir *Circuit) addInternalWire32FromSingleWires(singleWires [32]*wire.Wire) *wire.Wire32 {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	wire := wire.FromWires(wireName, singleWires)
	cir.internalWireCount++
	return wire
}

func (cir *Circuit) addComponents(components ...Component) {
	cir.components = append(cir.components, components...)
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

	// don't need to check if wire is in circuit;
	// if it isn't, the change won't propagate to anything in the circuit

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
func (cir *Circuit) Propagate(initialChanges ...Change) {
	agenda := NewChangeQueue()

	for _, ch := range initialChanges {
		agenda.AddChange(ch)
	}

	// fmt.Printf("Propagating from change at t=%v\n", initialChange.Time)

	for agenda.Length() > 0 {
		nextChanges, _ := agenda.GetNextChanges()

		// TODO - how to handle the case where there is >1 nextChange at the same time?
		// TODO - check for possibility of overlap and resolve/combine somehow?
		for _, nextChange := range nextChanges {
			// don't check if prevSignal == newSignal;
			// that's only an early-out optimization,
			// and not having it messes with circuit initialization
			downstreamChanges := cir.executeChange(nextChange)
			for _, downstream := range downstreamChanges {
				agenda.AddChange(downstream)
			}
		}
	}

	fmt.Printf("Done propagating\n\n")
}

func (cir *Circuit) Initialize() {
	inputInitializations := []Change{}

	for inWire := range cir.inputWires {
		inputInitializations = append(inputInitializations, Change{
			Time:   0,
			Wire:   inWire,
			Signal: wire.Low,
		})
	}

	for _, ch := range inputInitializations {
		cir.Propagate(ch)
	}
}
