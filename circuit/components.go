package circuit

import (
	"github.com/DylanSp/go-logical-circuits/circuit/wire"
)

func (cir *Circuit) MkNotGate(inWire *wire.Wire, outWire *wire.Wire) {
	gate := Component{
		inputWires: map[*wire.Wire]struct{}{},
	}
	gate.inputWires[inWire] = struct{}{}

	gate.processChangeLogic = func(ch Change) []Change {
		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   outWire,
				Signal: !ch.Signal,
			},
		}
	}

	cir.addComponents(gate)
}

func binaryGateFactory(truthTable func(wire.Signal, wire.Signal) wire.Signal) func(*Circuit, *wire.Wire, *wire.Wire, *wire.Wire) {
	return func(cir *Circuit, in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
		gate := Component{
			inputWires: map[*wire.Wire]struct{}{},
		}
		gate.inputWires[in1] = struct{}{}
		gate.inputWires[in2] = struct{}{}

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
	}
}

func (cir *Circuit) MkAndGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
	binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 && s2
	})(cir, in1, in2, out)
}

func (cir *Circuit) MkOrGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
	binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 || s2
	})(cir, in1, in2, out)
}

func (cir *Circuit) MkXorGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
	binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 != s2
	})(cir, in1, in2, out)
}

func (cir *Circuit) MkNandGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
	binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return !(s1 && s2)
	})(cir, in1, in2, out)
}

func (cir *Circuit) MkNorGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
	binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return !(s1 || s2)
	})(cir, in1, in2, out)
}

func (cir *Circuit) MkXnorGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
	binaryGateFactory(func(s1, s2 wire.Signal) wire.Signal {
		return s1 == s2
	})(cir, in1, in2, out)
}

func (cir *Circuit) Mk3AndGate(in1 *wire.Wire, in2 *wire.Wire, in3 *wire.Wire, out *wire.Wire) {
	tmp1 := cir.addInternalWire()
	cir.MkAndGate(in1, in2, tmp1)
	cir.MkAndGate(tmp1, in3, out)
}

func (cir *Circuit) MkHalfAdder(in1 *wire.Wire, in2 *wire.Wire, outSum *wire.Wire, outCarry *wire.Wire) {
	cir.MkXorGate(in1, in2, outSum)
	cir.MkAndGate(in1, in2, outCarry)
}
