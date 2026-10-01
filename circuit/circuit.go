package circuit

import "fmt"

// single change in a single wire
type Change struct {
	Time int

	Wire   *Wire
	Signal Signal
}

func (ch Change) String() string {
	return fmt.Sprintf("%v (t=%v) new value: %v", ch.Wire, ch.Time, ch.Signal)
}

// create a set of changes for a Wire32 based on `value`
func Change32(time int, wire *Wire32, value uint32) []Change {
	changes := make([]Change, 32)

	for i := range 32 {
		ch := Change{
			Time: time,
			Wire: wire.wires[i],
		}

		// TODO - based on isBitSet() from circuit_test; pull that out into a package for bit-level utility functions?
		bitmask := uint32(1) << i
		if (value & bitmask) != 0 {
			ch.Signal = High
		} else {
			ch.Signal = Low
		}

		changes[i] = ch
	}

	return changes
}

// processing delay of a single gate (of any type)
const gateDelay = 1

// Circuit defines the layout/connections of a logical circuit
type Circuit struct {
	// TODO - instead of tracking inputWires separately, initialize by setting a single wire (or all wires) to Low, then propagating?
	// TODO - would allow just having a single AddWire() method instead of Add{Input,Output,Auxiliary}Wire
	inputWires map[*Wire]struct{} // tracked for initialization
	wires      map[string]*Wire   // all wires by name; tracked to avoid duplicate wires

	internalWireCount int // used for giving internal wires unique names
}

func NewCircuit() Circuit {
	return Circuit{
		wires:             map[string]*Wire{},
		inputWires:        map[*Wire]struct{}{},
		internalWireCount: 0,
	}
}

// return the added wire so callers can refer to it for setting up components
func (cir *Circuit) addWire(wireName string) *Wire {
	_, hasWire := cir.wires[wireName]
	if hasWire {
		panic(fmt.Sprintf("Wire %v already present", wireName))
	}

	newWire := NewWire(wireName)
	cir.wires[wireName] = newWire

	return newWire
}

func (cir *Circuit) addWire32(wireName string) *Wire32 {
	// TODO - check if wire is already present? (circuit would need another field to track Wire32's)

	newWires := [32]*Wire{}
	for i := range 32 {
		newWires[i] = cir.addWire(fmt.Sprintf("%v-%v", wireName, i))
	}
	return FromWires(wireName, newWires)
}

func (cir *Circuit) AddInputWire(wireName string) *Wire {
	wire := cir.addWire(wireName)
	cir.inputWires[wire] = struct{}{}
	return wire
}

func (cir *Circuit) AddInputWire32(wireName string) *Wire32 {
	wire := cir.addWire32(wireName)

	// add the single-bit wires as input wires to the circuit
	for i := range 32 {
		cir.inputWires[wire.Wire(i)] = struct{}{}
	}

	return wire
}

// used for adding internal wires inside compound components
func (cir *Circuit) addInternalWire() *Wire {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	wire := cir.addWire(wireName)
	cir.internalWireCount++
	return wire
}

// TODO - do we need this?
func (cir *Circuit) addInternalWire32() *Wire32 {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	wire := cir.addWire32(wireName)
	cir.internalWireCount++
	return wire
}

func (cir *Circuit) addInternalWire32FromSingleWires(singleWires [32]*Wire) *Wire32 {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	wire := FromWires(wireName, singleWires)
	cir.internalWireCount++
	return wire
}
