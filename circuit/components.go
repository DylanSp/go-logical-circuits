package circuit

import (
	"fmt"

	"github.com/DylanSp/go-logical-circuits/circuit/wire"
)

func (cir *Circuit) Not(in *wire.Wire) *wire.Wire {
	gate := &Component{
		// inputWires: map[wire.Wire]struct{}{},
		inputWires: map[string]struct{}{},
	}
	gate.inputWires[in.Name()] = struct{}{}

	out := cir.addInternalWire()

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   out,
				Signal: !ch.Signal,
			},
		}
	}

	cir.addComponents(gate)

	return out
}

func binaryGateFactory(truthTable func(wire.Signal, wire.Signal) wire.Signal) func(*Circuit, *wire.Wire, *wire.Wire) *wire.Wire {
	return func(cir *Circuit, in1, in2 *wire.Wire) *wire.Wire {
		gate := &Component{
			// inputWires: map[wire.Wire]struct{}{},
			inputWires: map[string]struct{}{},
		}
		gate.inputWires[in1.Name()] = struct{}{}
		gate.inputWires[in2.Name()] = struct{}{}

		out := cir.addInternalWire()

		gate.processChangeLogic = func(ch Change) []Change {
			return []Change{
				{
					Time:   ch.Time + gateDelay,
					Wire:   out,
					Signal: truthTable(in1.Signal(), in2.Signal()),
				},
			}
		}
		cir.addComponents(gate)

		return out
	}
}

func (cir *Circuit) And(in1, in2 *wire.Wire) *wire.Wire {
	return binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 && s2
	})(cir, in1, in2)
}

func (cir *Circuit) Or(in1, in2 *wire.Wire) *wire.Wire {
	return binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 || s2
	})(cir, in1, in2)
}

func (cir *Circuit) Xor(in1, in2 *wire.Wire) *wire.Wire {
	return binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 != s2
	})(cir, in1, in2)
}

func (cir *Circuit) Nand(in1, in2 *wire.Wire) *wire.Wire {
	return binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return !(s1 && s2)
	})(cir, in1, in2)
}

func (cir *Circuit) Nor(in1, in2 *wire.Wire) *wire.Wire {
	return binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return !(s1 || s2)
	})(cir, in1, in2)
}

func (cir *Circuit) Xnor(in1, in2 *wire.Wire) *wire.Wire {
	return binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 == s2
	})(cir, in1, in2)
}

func (cir *Circuit) HalfAdder(in1, in2 *wire.Wire) (*wire.Wire, *wire.Wire) {
	sum := cir.Xor(in1, in2)
	carry := cir.And(in1, in2)

	return sum, carry
}

func (cir *Circuit) FullAdder(in1, in2, inCarry *wire.Wire) (*wire.Wire, *wire.Wire) {
	sum := cir.Xor(cir.Xor(in1, in2), inCarry)
	carry := cir.Or(cir.And(in1, in2), cir.And(inCarry, cir.Xor(in1, in2)))
	return sum, carry
}

// TODO - should this style of returning a component and the output wire be used for other gates/components?
func (cir *Circuit) Noop(in *wire.Wire) (*Component, *wire.Wire) {
	gate := &Component{
		// inputWires: map[wire.Wire]struct{}{},
		inputWires: map[string]struct{}{},
	}
	gate.inputWires[in.Name()] = struct{}{}

	out := cir.addInternalWire()

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   out,
				Signal: ch.Signal,
			},
		}
	}

	cir.addComponents(gate)

	return gate, out
}

func (cir *Circuit) Clock(period int) *wire.Wire {
	if period <= 0 {
		panic(fmt.Sprintf("Unable to create clock with period %v; period must be at least 1", period))
	}

	// delays := []Component{}
	// dummyWire := cir.addInternalWire()

	in := cir.addInternalWire() // dummy wire for initial input; will be irrelevant once we set up loopback
	var out *wire.Wire

	var initialNoop *Component

	// chain of (period - 1) no-ops to delay clock signal, followed by NOT gate to invert signal after `period` ticks
	for i := range period - 1 {
		gate, out := cir.Noop(in)
		if i == 0 {
			initialNoop = gate
		}

		in = out
	}

	// set up loopback, connect to first Noop component
	loopbackWire := cir.Not(out)
	initialNoop.AddInputWire(loopbackWire.Name())

	return out
}
