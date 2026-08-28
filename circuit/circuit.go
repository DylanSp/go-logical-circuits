package circuit

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit/wire"
)

// single change in a single wire
type Change struct {
	Time   int
	Wire   wire.Wire
	Signal bool
}

func (ch Change) String() string {
	return fmt.Sprintf("%v (t=%v) new value: %v", ch.Wire, ch.Time, ch.Signal)
}

// processing delay of a single gate (of any time)
const gateDelay = 2

// placeholder type;
// if I follow Haskell model, Components are functions that take a Change and a WireState, returning downstream changes
type Component struct {
	inputWires []*wire.Wire
}

func (c Component) ProcessChange(ch Change) {

}

type Circuit struct {
	wireState  map[string]*wire.Wire // wires by name
	components []Component

	internalWireCount int
}

func (cir *Circuit) AddComponents(components ...Component) {
	panic("unimplemented")
}

func (cir *Circuit) AddWire(wireName string) {
	_, hasWire := cir.wireState[wireName]
	if hasWire {
		panic(fmt.Sprintf("Wire %v already present", wireName))
	}

	newWire := wire.New(wireName)
	cir.wireState[wireName] = &newWire
}

// used for adding internal wires inside components
// TODO - does this need to be exported?
func (cir *Circuit) AddInternalWire() {
	wireName := fmt.Sprintf("internal-%v", cir.internalWireCount)
	cir.AddWire(wireName)
	cir.internalWireCount++
}
