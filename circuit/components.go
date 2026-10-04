package circuit

import (
	"fmt"
	"maps"
	"slices"
)

func (cir *Circuit) Not(in *Wire) *Wire {
	gate := cir.addInternalWire()
	gate.SetInputs(in)

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   gate,
				Signal: !ch.Signal,
			},
		}
	}

	return gate
}

func binaryGateFactory(truthTable func(Signal, Signal) Signal) func(*Circuit, *Wire, *Wire) *Wire {
	return func(cir *Circuit, in1, in2 *Wire) *Wire {
		gate := cir.addInternalWire()
		gate.SetInputs(in1, in2)

		gate.processChangeLogic = func(ch Change) []Change {
			currentInputs := slices.Collect(maps.Values(gate.inputWires))

			return []Change{
				{
					Time:   ch.Time + gateDelay,
					Wire:   gate,
					Signal: truthTable(currentInputs[0].Signal(), currentInputs[1].Signal()),
				},
			}
		}

		return gate
	}
}

func (cir *Circuit) And(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 && s2
	})(cir, in1, in2)
}

func (cir *Circuit) Or(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 || s2
	})(cir, in1, in2)
}

func (cir *Circuit) Xor(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 != s2
	})(cir, in1, in2)
}

func (cir *Circuit) Nand(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return !(s1 && s2)
	})(cir, in1, in2)
}

func (cir *Circuit) Nor(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return !(s1 || s2)
	})(cir, in1, in2)
}

func (cir *Circuit) Xnor(in1, in2 *Wire) *Wire {
	return binaryGateFactory(func(s1, s2 Signal) Signal {
		return s1 == s2
	})(cir, in1, in2)
}

func (cir *Circuit) HalfAdder(in1, in2 *Wire) (*Wire, *Wire) {
	sum := cir.Xor(in1, in2)
	carry := cir.And(in1, in2)

	return sum, carry
}

func (cir *Circuit) FullAdder(in1, in2, inCarry *Wire) (*Wire, *Wire) {
	sum := cir.Xor(cir.Xor(in1, in2), inCarry)
	carry := cir.Or(cir.And(in1, in2), cir.And(inCarry, cir.Xor(in1, in2)))
	return sum, carry
}

func (cir *Circuit) Noop(in *Wire) *Wire {
	gate := cir.addInternalWire()
	gate.SetInputs(in)

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   gate,
				Signal: ch.Signal,
			},
		}
	}

	return gate
}

// TODO - doesn't currently work if period=1; try and fix this?
func (cir *Circuit) Clock(period int) *Wire {
	if period <= 0 {
		panic(fmt.Sprintf("Unable to create clock with period %v; period must be at least 1", period))
	}

	in := cir.addInternalWire() // dummy wire for initial input when constructing no-ops; will be irrelevant once we set up loopback
	var out *Wire

	var initialNoop *Wire

	// chain of (period - 1) no-ops to delay clock signal, followed by NOT gate to invert signal after `period` ticks
	for i := range period - 1 {
		out = cir.Noop(in)
		if i == 0 {
			initialNoop = out
		}

		in = out
	}

	// set up loopback, connect to first Noop component
	loopbackWire := cir.Not(out)
	initialNoop.SetInputs(loopbackWire)

	return out
}
