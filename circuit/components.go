package circuit

import (
	"github.com/DylanSp/go-logical-circuits/circuit/wire"
)

func (cir *Circuit) Not(in *wire.Wire) *wire.Wire {
	gate := Component{
		inputWires: map[*wire.Wire]struct{}{},
	}
	gate.inputWires[in] = struct{}{}

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
		gate := Component{
			inputWires: map[*wire.Wire]struct{}{},
		}
		gate.inputWires[in1] = struct{}{}
		gate.inputWires[in2] = struct{}{}

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
