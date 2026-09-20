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
	fmt.Printf("Calling SetSignal on wire %v with value %v\n", w, newSignal)
	w.signal = newSignal
}
