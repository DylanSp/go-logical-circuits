package wire

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
}

func New(name string) *Wire {
	return &Wire{
		name:   name,
		signal: Low,
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

func (w *Wire) SetSignal(newSignal Signal) {
	// fmt.Printf("Calling SetSignal on wire %v with value %v\n", w, newSignal)
	w.signal = newSignal
}

type Wire32 struct {
	name  string
	wires [32]*Wire
}

func New32(name string) *Wire32 {
	newWire32 := &Wire32{
		name:  name,
		wires: [32]*Wire{},
	}

	for i := range newWire32.wires {
		newWire := New(fmt.Sprintf("%v-%v", name, i))
		newWire32.wires[i] = newWire
	}

	return newWire32
}

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
