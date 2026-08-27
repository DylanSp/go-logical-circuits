package circuit

import "fmt"

type Wire struct {
	name string
}

func (w Wire) Name() string {
	return w.name
}

// implement stringer interface
func (w Wire) String() string {
	return w.Name()
}

// signals will be represented as booleans for now

// single change in a single wire
type Change struct {
	Time   int
	Wire   Wire
	Signal bool
}

func (ch Change) String() string {
	return fmt.Sprintf("%v (t=%v) new value: %v", ch.Wire, ch.Time, ch.Signal)
}

// placeholder type;
// if I follow Haskell model, Components are functions that take a Change and a WireState, returning downstream changes
type Component int

// TODO - tracking wire state
// should Wires track their own state, then Circuit.wireState can just have a set of wires?
type Circuit struct {
	wireState  map[Wire]bool
	components []Component

	internalWireCount int
}

func (cir *Circuit) AddComponents(components ...Component) {
	panic("unimplemented")
}

func (cir *Circuit) AddWire(wireName string) {
	newWire := Wire{
		name: wireName,
	}

	_, ok := cir.wireState[newWire]
	if ok {
		panic(fmt.Sprintf("Wire %v already present", wireName))
	} else {
		cir.wireState[newWire] = false
	}
}

// used for adding internal wires inside components
// TODO - does this need to be exported?
func (cir *Circuit) AddInternalWire() {
	newWire := Wire{
		name: fmt.Sprintf("internal-%v", cir.internalWireCount),
	}
	cir.wireState[newWire] = false

	cir.internalWireCount++
}
