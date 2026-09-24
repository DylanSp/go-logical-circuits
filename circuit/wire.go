package circuit

import "fmt"

type Signal bool

const (
	Low  Signal = false
	High Signal = true
)

func (sig Signal) String() string {
	if sig == Low {
		return "Low"
	} else {
		return "High"
	}
}

type Wire struct {
	name   string
	signal Signal

	// track input wires by name
	inputWires map[string]struct{}

	// internal logic for responding to changes
	// does *not* update wires' states; Circuit.ExecuteChanges() will take the changes from this and update states appropriately
	processChangeLogic func(Change) []Change
}

func NewWire(name string) *Wire {
	return &Wire{
		name:       name,
		signal:     Low,
		inputWires: map[string]struct{}{},
	}
}

func (w *Wire) Name() string {
	return w.name
}

// implement stringer interface
func (w *Wire) String() string {
	return w.Name()
}

func (w *Wire) Signal() Signal {
	return w.signal
}

// utility method to make testing easier
func (w *Wire) IsHigh() bool {
	return bool(w.Signal())
}

func (w *Wire) SetSignal(newSignal Signal) {
	// fmt.Printf("Calling SetSignal on wire %v with value %v\n", w, newSignal)
	w.signal = newSignal
}

// TODO - is this needed?
// adds an input wire to existing inputs
func (w *Wire) AddInput(newInput *Wire) {
	w.inputWires[newInput.Name()] = struct{}{}
}

// TODO - is this needed?
// resets input wires to the new inputs
func (w *Wire) SetInputs(newInputs ...*Wire) {
	w.inputWires = map[string]struct{}{}

	for _, in := range newInputs {
		w.inputWires[in.Name()] = struct{}{}
	}
}

func (w *Wire) HandleChange(ch Change) []Change {
	// fmt.Println("Input wires for component:")
	// for inputWire := range w.inputWires {
	// fmt.Printf("Wire %v\n", inputWire)

	// fmt.Printf("Wire %v\n", inputWire.Name())
	// fmt.Printf("  Memory address: %p\n", inputWire)
	// }
	// fmt.Println()

	// fmt.Println("Wire from change:")
	// fmt.Printf("Wire %v\n", ch.Wire.Name())
	// fmt.Printf("  Memory address: %p\n", ch.Wire)
	// fmt.Println()

	if _, ok := w.inputWires[ch.Wire.Name()]; !ok {
		// fmt.Println("Wire is not an input wire of this component; skipping")

		// the change isn't one of this component's input wires; no changes produced
		return []Change{}
	}

	// fmt.Println("Processing change")

	return w.processChangeLogic(ch)
}

type Wire32 struct {
	name  string
	wires [32]*Wire
}

func NewWire32(name string) *Wire32 {
	newWire32 := &Wire32{
		name:  name,
		wires: [32]*Wire{},
	}

	for i := range newWire32.wires {
		newWire := NewWire(fmt.Sprintf("%v-%v", name, i))
		newWire32.wires[i] = newWire
	}

	return newWire32
}

// TODO - rename to be less ambiguous, since this is now in the circuit package?
func FromWires(name string, singleWires [32]*Wire) *Wire32 {
	newWire32 := &Wire32{
		name:  name,
		wires: singleWires,
	}
	return newWire32
}

func (w *Wire32) Name() string {
	return w.name
}

// implement stringer interface
func (w *Wire32) String() string {
	return w.Name()
}

func (w *Wire32) AllWires() [32]*Wire {
	return w.wires
}

func (w *Wire32) Wire(i int) *Wire {
	return w.wires[i]
}

func (w *Wire32) AsUint32() uint32 {
	value := uint32(0)

	for i := range 32 {
		if w.Wire(i).Signal() == High {
			// set i'th bit of value
			value |= (uint32(1) << i)
		}
	}

	return value
}
