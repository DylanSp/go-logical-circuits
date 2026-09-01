package circuit

import (
	"fmt"

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

func (cir *Circuit) MkAndGate(in1 *wire.Wire, in2 *wire.Wire, out *wire.Wire) {
	gate := Component{
		inputWires: map[*wire.Wire]struct{}{},
	}
	gate.inputWires[in1] = struct{}{}
	gate.inputWires[in2] = struct{}{}

	gate.processChangeLogic = func(ch Change) []Change {
		fmt.Printf("AND gate processing change at t=%v\n", ch.Time)
		fmt.Printf("Change on wire %v, new signal is %v\n", ch.Wire, ch.Signal)

		return []Change{
			{
				Time:   ch.Time + gateDelay,
				Wire:   out,
				Signal: in1.Signal() && in2.Signal(),
			},
		}
	}

	cir.addComponents(gate)
}

func (cir *Circuit) Mk3AndGate(in1 *wire.Wire, in2 *wire.Wire, in3 *wire.Wire, out *wire.Wire) {
	tmp1 := cir.addInternalWire()
	cir.MkAndGate(in1, in2, tmp1)
	cir.MkAndGate(tmp1, in3, out)
}
