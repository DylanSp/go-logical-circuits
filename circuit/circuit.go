package circuit

import "fmt"

type Wire struct {
	name string
}

func (w Wire) Name() string {
	return w.name
}

// signals will be represented as booleans for now

type Component int // placeholder

type CircuitState struct {
	wireState  map[Wire]bool
	components []Component
}

// single change in a single wire
type Change struct {
	Time   int
	Wire   Wire
	Signal bool
}

func (ch Change) String() string {
	return fmt.Sprintf("%v (t=%v) new value: %v", ch.Wire.Name(), ch.Time, ch.Signal)
}
